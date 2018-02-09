package main
import(
  "github.com/boltdb/bolt"
  "net/http"
  "fmt"
  "html/template"
  "time"
  "encoding/gob"
  "bytes"
)

type PageVariables struct{
  Title string
}

type User struct{
  Username string
  Password string
  CreatedOn time.Time
}

func Update(key string, user *User) error {
  var res bytes.Buffer
  encoder := gob.NewEncoder(&res)
  encoder.Encode(user)
  encodedRecord := res.Bytes()
	err := db.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte("USER"))
		if e != nil {
			return e
		}
		// Encode the record
		e = bucket.Put([]byte(key), encodedRecord)
     if e != nil {
			return e
		}
		return nil
	})
	return err
}


func View(key string, pwd string) error{
  err := db.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte("USER"))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		// Retrieve the record
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}else{
      fmt.Println("getuser")
      var user User
      decoder := gob.NewDecoder(bytes.NewReader(v))
      decoder.Decode(&user)
      fmt.Println(&user)
      dbpwd := user.Password
      fmt.Println(dbpwd)
    }
    return nil
  })
  return err
}

func createHandler(w http.ResponseWriter, r *http.Request){
  title := PageVariables{"Create a New Account"}
  t, err := template.ParseFiles("newaccount.html")
  if err == nil{
    t.Execute(w, title)
  }else{
    fmt.Println(err)
  }
}

func loginHandler(w http.ResponseWriter, r *http.Request){
  title := PageVariables{"Log In"}
  t, err := template.ParseFiles("login.html")
  if err == nil{
    t.Execute(w, title)
  }else{
    fmt.Println(err)
  }
}

func signinHandler(w http.ResponseWriter, r *http.Request){
  fmt.Println("hello")
  username := r.FormValue("username")
  password := r.FormValue("password")
  View(username, password)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
  username := r.FormValue("username")
  password := r.FormValue("password")
  user := User{username, password, time.Now()}
  Update(username, &user)
}

var db *bolt.DB

func main(){
  var err error
  db,err = bolt.Open("demo.db", 0600, nil)
  if err != nil{
    fmt.Println(err)
  }
  http.HandleFunc("/signon/", createHandler)
  http.HandleFunc("/createaccount/",saveHandler)
  http.HandleFunc("/login/",loginHandler)
  http.HandleFunc("/signinaccount/", signinHandler)
  http.ListenAndServe(":8080",nil)
}
