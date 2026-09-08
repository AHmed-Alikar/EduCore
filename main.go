func studentHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	// /students
	if len(parts) == 2 {

		switch r.Method {

		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(students)

		case http.MethodPost:
			var student Student

			err := json.NewDecoder(r.Body).Decode(&student)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid JSON")
				return
			}

			student.ID = 1003
			students[student.ID] = student

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(student)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, "Method Not Allowed")
		}

		return
	}

	// /students/{id}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid Student ID")
		return
	}

	switch r.Method {

	case http.MethodGet:
		student, exists := students[id]

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Student not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(student)

	case http.MethodPut:
		var student Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid JSON")
			return
		}

		_, exists := students[id]

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Student not found")
			return
		}

		student.ID = id
		students[id] = student

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(student)

	case http.MethodDelete:
		_, exists := students[id]

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Student not found")
			return
		}

		delete(students, id)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Student deleted successfully")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method Not Allowed")
	}
}