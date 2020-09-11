package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

// NumbersHandler function
func NumbersHandler(ctx echo.Context) error {
	number := ctx.Param("number")

	url := "http://numbersapi.com/" + number

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("invalid number parameter: %s", number)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return ctx.JSON(http.StatusOK, string(body))
}

Step 3: Create the unit test api/api_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func TestNumbersHandler(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	type row struct {
		name    string
		args    args
		wantErr bool
	}
	tests := make([]row, 0, 2)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/numbersapi/:number")
	c.SetParamNames("number")
	c.SetParamValues("1")
	tests = append(tests, row{
		name:    "numberhandler",
		args:    args{c},
		wantErr: false,
	})

	d := e.NewContext(req, rec)
	d.SetPath("/numbersapi/:number")
	d.SetParamNames("number")
	d.SetParamValues("notanumber")
	tests = append(tests, row{
		name:    "numberhandler",
		args:    args{d},
		wantErr: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NumbersHandler(tt.args.ctx); (err != nil) != (tt.wantErr == true) {
				t.Errorf("NumbersHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}