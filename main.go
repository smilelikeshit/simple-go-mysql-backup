//https://github.com/JamesStewy/go-mysqldump

package main

import (
	"fmt"
	"database/sql"
	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/exec"
	"log"
	"github.com/joho/godotenv"
	"github.com/jasonlvhit/gocron"

)

const (
	PATH_DAILY = "dumps/daily"
	PATH_WEEKLY = "dumps/weekly"
	PATH_MONTHLY = "dumps/monthly"
)

func initDB() *sql.DB{

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME))
	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}
	return db

}

func Daily(){

	db := initDB()

	DB_NAME := os.Getenv("DB_NAME")

	_, err := os.Stat(PATH_DAILY)
 
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(PATH_DAILY, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
 
	}

	dumpFilenameFormat := fmt.Sprintf("%s-daily-20060102T150405", DB_NAME)   // accepts time layout string and add .sql at the end of file

	dumper, err := mysqldump.Register(db, PATH_DAILY, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	fmt.Printf("\nFile is saved to %s", resultFilename)

	// Close dumper and connected database
	dumper.Close()

}

func Weekly(){

	db := initDB()

	DB_NAME := os.Getenv("DB_NAME")

	_, err := os.Stat(PATH_WEEKLY)
 
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(PATH_WEEKLY, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
 
	}

	dumpFilenameFormat := fmt.Sprintf("%s-weekly-20060102T150405", DB_NAME)   // accepts time layout string and add .sql at the end of file

	dumper, err := mysqldump.Register(db, PATH_WEEKLY, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	fmt.Printf("\nFile is saved to %s", resultFilename)

	// Close dumper and connected database
	dumper.Close()

}

func RotateDaily(){

	arg1 := "find dumps/daily -type f -mtime +30 -exec rm -rf {};"
	cmd := exec.Command(arg1)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
	}
	fmt.Println(stdout)
}

func RotateWeekly(){

	arg1 := "find dumps/weekly -type f -mtime +120 -exec rm -rf {};"
	cmd := exec.Command(arg1)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
	}
	fmt.Println(stdout)
}

func Echo(){
	fmt.Println("service up !")
}


func main(){

	fmt.Println("Service up")
	s := gocron.NewScheduler()

	s.Every(30).Second().Do(Echo)

	s.Every(1).Day().Do(Daily)

	s.Every(1).Week().Do(Weekly)

	// clean files old 1 months
	s.Every(1).Day().Do(RotateDaily)

	// clean files old 3 months
	s.Every(1).Day().Do(RotateWeekly)

	
	<- s.Start()
}


