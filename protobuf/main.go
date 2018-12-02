package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"golearn/protobuf/document"
	"io"
	"os"
	"time"
)

func encodeInt(x int) []byte {
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, uint64(x))
	return buf[0:n]
}


func writeDocs(writer io.Writer) error {
	urls := []string{
		"https://roem.ru/20-11-2018/275185/podderzhat-neponyatno/",
		"https://news.mail.ru/politics/35561378",
		"https://developers.google.com/protocol-buffers/docs/gotutorial",
		"https://www.techradar.com/reviews/anmeldelse-ipad-pro-11",
	}

	for _, u := range urls {
		tm := time.Now()
		doc := &document.Document{
			Url: u,
			Body: fmt.Sprintf("Content of %s\n", u),
			DownloadTime: &timestamp.Timestamp{Seconds: tm.Unix(), Nanos: int32(tm.UnixNano())},
		}

		data, err := proto.Marshal(doc)
		if err != nil {
			return errors.Wrap(err, "failed to marshal proto")
		}

		if _, err = writer.Write(encodeInt(len(data))); err == nil {
			_, err = writer.Write(data)
		}
		if err != nil {
			return errors.Wrap(err, "failed to write to file")
		}
	}

	return nil
}

func readDocs(reader io.Reader) ([]document.Document, error) {
	var docs []document.Document
	rd := bufio.NewReader(reader)

	for {
		size, err := binary.ReadUvarint(rd)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "failed to decode size")
		}

		buffer := make([]byte, size)
		n, err := rd.Read(buffer)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read file")
		}
		if n < int(size) {
			return nil, errors.Wrap(err, "unexpected EOF")
		}

		var doc document.Document
		if err := proto.Unmarshal(buffer, &doc); err != nil {
			return nil, errors.Wrap(err, "failed to decode message")
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

func protoTimeToUnix(tm *timestamp.Timestamp) time.Time {
	return time.Unix(tm.Seconds, int64(tm.Nanos))
}

func printDocs(reader io.Reader) error {
	docs, err := readDocs(reader)
	if err != nil {
		return errors.Wrap(err, "failed to load documents")
	}

	for _, doc := range docs {
		fmt.Printf("%s:\n", doc.GetUrl())
		fmt.Printf("\tBody: %d bytes\n", len(doc.GetBody()))
		fmt.Printf("\tDownloaded: %s\n", protoTimeToUnix(doc.GetDownloadTime()).UTC())
		fmt.Println()
	}

	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v", err)
		os.Exit(1)
	}
}

func main() {
	usage := func() {
		fmt.Fprintln(os.Stderr, "Usage: protobuf {read|write} file")
		os.Exit(2)
	}

	if len(os.Args) != 3 {
		usage()
	}

	cmd, file := os.Args[1], os.Args[2]

	switch cmd {
	case "write":
		fd, err := os.OpenFile(file, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
		checkErr(err)
		checkErr(writeDocs(fd))
	case "read":
		fd, err := os.Open(file)
		checkErr(err)
		checkErr(printDocs(fd))
	default:
		usage()
	}
}
