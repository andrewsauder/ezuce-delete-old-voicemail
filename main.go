package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type voicemail_file struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Filename   string
	UploadDate time.Time `bson:"uploadDate"`
}
type voicemail_chunk struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	files_id bson.ObjectId `bson:"files_id,omitempty"`
}

func main() {
	//define required arguements
	flag.Usage = func() {
		fmt.Println("-------------------------------------------------------------------------------")
		fmt.Println("    DELOLDVM     ::     Delete Ezuce UniteMe voicemail older than X months")
		fmt.Println("-------------------------------------------------------------------------------")
		flag.PrintDefaults()
	}

	var url string
	var username string
	var password string
	var months int
	flag.StringVar(&url, "url", "127.0.0.1:27017", "Specify a mongo connection url to use including port")
	flag.StringVar(&username, "u", "", "Specify Mongodb username. Leave blank for no authentication.")
	flag.StringVar(&password, "p", "", "Specify Mongodb password. Leave blank for no authentication.")
	flag.IntVar(&months, "months", 6, "Specify the number of months to keep voicemail for (ie. entering 6 will delete all voicemails older than 6 months)")
	flag.Parse()

	//set up log
	//create your file with desired read/write permissions
	f, err := os.OpenFile("delete-old-voicemail.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set the log output
	log.SetOutput(f)
	log.Println("Running")

	//mongo connection details
	//url := "config.uc.garrettcounty.org:27017"
	database := "vmdb"
	collection_chunks_name := "fs.chunks"
	collection_files_name := "fs.files"

	//connect to mongo
	info := &mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	//get collections
	c_files := session.DB(database).C(collection_files_name)
	c_chunks := session.DB(database).C(collection_chunks_name)

	//set date to delete older than
	deleteOlderThan := time.Now().AddDate(0, months*-1, 0)

	//query for files with older upload times
	results := []voicemail_file{}
	err = c_files.Find(
		bson.M{
			"metadata.user": "5001",
			"uploadDate": bson.M{
				"$lt": deleteOlderThan,
			},
		}).Sort("-uploadDate").All(&results)

	if err != nil {
		panic(err)
	}

	//loop the files that we're going to delete
	for _, file := range results {
		log.Println("File: " + file.ID.Hex())
		log.Println("    Uploaded: " + file.UploadDate.String())

		//get the chunks associated with the file
		chunks := []voicemail_chunk{}
		err = c_chunks.Find(bson.M{"files_id": file.ID}).All(&chunks)
		if err != nil {
			panic(err)
		}

		//loop the chunks that we need to delete before we delete the file
		for _, chunk := range chunks {
			log.Println("     CHUNK: " + chunk.ID.Hex())
			//TODO: delete chunk
			//c_chunks.RemoveId(chunk.ID)
		}

		//TODO: delete message
		//c_files.RemoveId(file.ID)
	}
}
