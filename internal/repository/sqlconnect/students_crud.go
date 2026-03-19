package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/pkg/utils"
	"strconv"
	"strings"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func GetStudentsDbHandler(students []models.Student, r *http.Request) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connecting to DataBase")
	}
	defer db.Close()

	//teacher
	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	query = utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "Database query error")
	}
	defer rows.Close()

	// teacherList := make([]models.Teacher, 0)
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error retrieving data")
		}
		students = append(students, student)
	}
	return students, nil
}

func GetStudentByID(id int) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error connecting to database")
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "Error retrieving data")
	} else if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error retrieving data")
	}
	return student, nil
}

func AddStudentsDbHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connecting to database")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", models.Student{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error adding data")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))
	for i, newStudent := range newStudents {
		// res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		values := utils.GetStructValues(newStudent)
		res, err := stmt.Exec(values...)
		if err != nil {
			fmt.Println("-------------Error:", err.Error())
			if strings.Contains(err.Error(), "a foreign key constraint fails (`school`.`students`, CONSTRAINT `students_ibfk_1` FOREIGN KEY (`class`) REFERENCES `yeachers` (`class`))") {
				return nil, utils.ErrorHandler(err, "class/class teacher dot not exist")
			}
			return nil, utils.ErrorHandler(err, "Error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error getting last insert ID")
		}
		newStudent.ID = int(lastID)
		addedStudents[i] = newStudent
	}
	return addedStudents, nil
}

func UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "Unable to connect to database")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?",
		id).Scan(&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName,
		&existingStudent.Email, &existingStudent.Class)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return models.Student{}, utils.ErrorHandler(err, "Student not found")
		}
		log.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	updatedStudent.ID = existingStudent.ID
	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", updatedStudent.FirstName,
		updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, updatedStudent.ID)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating data")
	}
	return updatedStudent, nil
}

func PatchStudents(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return utils.ErrorHandler(err, "Unable to connect to database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return utils.ErrorHandler(err, "Error starting transaction")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid student ID in update")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error converting ID to int")
		}

		var studentFromDb models.Student
		err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&studentFromDb.ID,
			&studentFromDb.FirstName, &studentFromDb.LastName, &studentFromDb.Email, &studentFromDb.Class)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				tx.Rollback()
				return utils.ErrorHandler(err, "Student not found")
			}
			return utils.ErrorHandler(err, "Error retrieving student")
		}

		// Apply updates using reflection
		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skip updating the id field
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("cannot convert %v to %v", val.Type(), fieldVal.Type())
							return utils.ErrorHandler(err, "Error updating data")
						}

					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", studentFromDb.FirstName,
			studentFromDb.LastName, studentFromDb.Email, studentFromDb.Class, studentFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating data")
		}

	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error committing transaction")
	}
	return nil
}

func PatchOneStudent(id int, updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "Unable to connect to database")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return models.Student{}, utils.ErrorHandler(err, "Error retrieving data")
		}
		log.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "Error retrieving data")
	}

	studentVal := reflect.ValueOf(&existingStudent).Elem()
	studentType := studentVal.Type()

	for k, v := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					fieldVal := studentVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", existingStudent.FirstName,
		existingStudent.LastName, existingStudent.Email, existingStudent.Class, existingStudent.ID)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating data")
	}
	return existingStudent, nil
}

func DeleteOneStudent(id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return utils.ErrorHandler(err, "unable to connect to database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting data")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error retrieving delete data")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Error retrieving data")
	}
	return nil
}

func DeleteStudents(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return nil, utils.ErrorHandler(err, "Unable to connect to database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, utils.ErrorHandler(err, "Error starting transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error deleting data")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return nil, utils.ErrorHandler(err, "Error deleting data")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error retrieving data")
		}

		// if teacher was deleted the add the ID to the deletedIds slice
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

		if rowsAffected < 1 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("ID %d not found", id))
		}

	}

	// Commit
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, utils.ErrorHandler(err, "Error committin transaction")
	}

	if len(deletedIds) < 1 {
		log.Println(err)
		return nil, utils.ErrorHandler(err, "Error deleting data")
	}
	return deletedIds, nil
}
