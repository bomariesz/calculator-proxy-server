package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func transformJSON(v interface{}) ([]byte, error) {
	dataJSON, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

type calcData struct {
	Result float64
}

type errorData struct {
	Message string
}

func transformErrorResponse(responseBody io.ReadCloser, w http.ResponseWriter) (errorData, error) {
	decoder := json.NewDecoder(responseBody)

	var errorResponse errorData
	err := decoder.Decode(&errorResponse)

	return errorResponse, err
}

func transformResponse(responseBody io.ReadCloser, w http.ResponseWriter) (calcData, error) {
	decoder := json.NewDecoder(responseBody)

	var response calcData
	err := decoder.Decode(&response)

	return response, err
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

func respServiceUnavailable(data errorData, w http.ResponseWriter) {
	dataJSON, err := transformJSON(data)
	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}

		respInternalError(errData, w)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write(dataJSON)
}

func respAPIOriginError(responseStatusCode int, responseBody io.ReadCloser, w http.ResponseWriter) {
	responseData, err := transformErrorResponse(responseBody, w)
	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}
		respInternalError(errData, w)
	} else {
		errorJSON, err := transformJSON(responseData)
		if err != nil {
			errData := errorData{
				Message: err.Error(),
			}

			respInternalError(errData, w)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(responseStatusCode)
		w.Write(errorJSON)
	}
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

func callCalculatorAPI(w http.ResponseWriter, r *http.Request, url string) {
	response, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		errData := errorData{
			Message: err.Error(),
		}
		respServiceUnavailable(errData, w)
	} else if response.StatusCode != http.StatusOK {
		respAPIOriginError(response.StatusCode, response.Body, w)
	} else {
		responseData, err := transformResponse(response.Body, w)
		if err != nil {
			errData := errorData{
				Message: err.Error(),
			}
			respInternalError(errData, w)
		} else {
			respSuccess(responseData, w)
		}
	}
}

func addition(w http.ResponseWriter, r *http.Request) {

}

func subtraction(w http.ResponseWriter, r *http.Request) {

}

func multiplication(w http.ResponseWriter, r *http.Request) {

}

func division(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/calculator/sum", addition)
	mux.HandleFunc("/calculator/sub", subtraction)
	mux.HandleFunc("/calculator/mul", multiplication)
	mux.HandleFunc("/calculator/div", division)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
