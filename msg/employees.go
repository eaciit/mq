package msg

import(
  "fmt"
  "encoding/json"
  "strconv"
)

// json structure in value interface{}
// {"name":"nanda","role":"petugas"}
type Employee struct {
	Name string
	Role string
}


func Marshal(value interface{}) string{
	js, err := json.Marshal(value)
    if err != nil {
        return ""
  }
  return string(js)
}


func GetEmployeeRole(value interface{}) string{
  var employee Employee
	strVal := Marshal(value)
  s, _:= strconv.Unquote(strVal)
	err := json.Unmarshal([]byte(s),&employee)
  if err != nil {
      fmt.Println(err)
  }
	return employee.Role
}
