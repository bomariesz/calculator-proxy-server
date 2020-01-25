package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type calcRequest struct {
	A float64
	B float64
}

type calcData struct {
	Result float64
}

type errorData struct {
	Message string
}

func transformJSON(v interface{}) ([]byte, error) {
	dataJSON, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func respSuccess(data calcData, w http.ResponseWriter) {
	dataJSON, err := transformJSON(data)
	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respInternalError(errData, w)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataJSON)
}

func respBadRequest(data errorData, w http.ResponseWriter) {
	dataJSON, err := transformJSON(data)
	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respInternalError(errData, w)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(dataJSON)
}

func respInternalError(data errorData, w http.ResponseWriter) {
	dataJSON, err := transformJSON(data)
	if err != nil {
		log.Panicln(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(dataJSON)
}

func transformRequest(w http.ResponseWriter, r *http.Request) (calcRequest, error) {
	decoder := json.NewDecoder(r.Body)

	var request calcRequest
	err := decoder.Decode(&request)

	return request, err
}

func calculate(w http.ResponseWriter, variable calcRequest, operation string) {
	result := 0.0
	var err error
	err = nil

	switch operation {
	case "+":
		result = variable.A + variable.B
	case "-":
		result = variable.A - variable.B
	case "*":
		result = variable.A * variable.B
	case "/":
		if variable.B == 0 {
			err = errors.New("you tried to divide by zero")
		} else {
			result = variable.A / variable.B
		}
	default:
		err = errors.New("invalid operation selected")
	}

	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respBadRequest(errData, w)
	} else {
		resultData := calcData{
			Result: result,
		}

		respSuccess(resultData, w)
	}
}

func addition(w http.ResponseWriter, r *http.Request) {
	reqCalc, err := transformRequest(w, r)

	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respBadRequest(errData, w)
	} else {
		calculate(w, reqCalc, "+")
	}
}

func subtraction(w http.ResponseWriter, r *http.Request) {
	reqCalc, err := transformRequest(w, r)

	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respBadRequest(errData, w)
	} else {
		calculate(w, reqCalc, "-")
	}
}

func multiplication(w http.ResponseWriter, r *http.Request) {
	reqCalc, err := transformRequest(w, r)

	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respBadRequest(errData, w)
	} else {
		calculate(w, reqCalc, "*")
	}
}

func division(w http.ResponseWriter, r *http.Request) {
	reqCalc, err := transformRequest(w, r)

	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respBadRequest(errData, w)
	} else {
		calculate(w, reqCalc, "/")
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/calculator.sum":
		if r.Method == http.MethodPost {
			addition(w, r)
		} else {
			http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		}
	case "/calculator.sub":
		if r.Method == http.MethodPost {
			subtraction(w, r)
		} else {
			http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		}
	case "/calculator.mul":
		if r.Method == http.MethodPost {
			multiplication(w, r)
		} else {
			http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		}
	case "/calculator.div":
		if r.Method == http.MethodPost {
			division(w, r)
		} else {
			http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", requestHandler)

	log.Fatal(http.ListenAndServe(":8090", mux))
}
