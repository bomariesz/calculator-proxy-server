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
			err = errors.New("error: you tried to divide by zero")
		} else {
			result = variable.A / variable.B
		}
	default:
		err = errors.New("error: Invalid operation selected")
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/calculator/sum", addition)
	mux.HandleFunc("/calculator/sub", subtraction)
	mux.HandleFunc("/calculator/mul", multiplication)
	mux.HandleFunc("/calculator/div", division)

	log.Fatal(http.ListenAndServe(":8090", mux))
}
