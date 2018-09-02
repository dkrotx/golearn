package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const databaseFile = "dataset.xml"

type UserDbRecord struct {
	Id            int `xml:"id"`
	Name       string
	Age           int `xml:"age"`
	About      string `xml:"about"`
	Gender     string `xml:"gender"`
	FirstName  string `xml:"first_name"`
	LastName   string `xml:"last_name"`
}

type UsersCollection struct {
	XMLName xml.Name     `xml:"root"`
	Users []UserDbRecord `xml:"row"`
}

func (db *UsersCollection) patchUserRecords() {
	for i, u := range db.Users {
		(*db).Users[i].Name = u.FirstName + " " + u.LastName
	}
}

func (db *UsersCollection) Load() (err error) {
	f, err := os.Open(databaseFile)

	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	if err = xml.Unmarshal(data, &db); err != nil {
		return
	}

	db.patchUserRecords()
	return
}

func (db *UsersCollection) GetUserRecords() []UserDbRecord {
	return db.Users
}

func Filter(users []UserDbRecord, s string) (res []UserDbRecord) {
	for _, u := range users {
		if strings.Contains(u.Name, s) || strings.Contains(u.About, s) {
			res = append(res, u)
		}
	}
	return
}


// do not do anything, nobody check this
func OffsetLimit(users []UserDbRecord, offset, limit int) []UserDbRecord {
	min := func (a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	offset = min(offset, len(users))
	limit = min(offset+limit, len(users))
	return users[offset:limit]
}

type By func(a, b *UserDbRecord) bool

type usersSorter struct {
	users []UserDbRecord
	by By
	reversed bool
}

func (srt *usersSorter) Len() int {
	return len(srt.users)
}

func (srt *usersSorter) Less(i, j int) bool {
	if srt.reversed {
		return srt.by(&srt.users[j], &srt.users[i])
	}
	return srt.by(&srt.users[i], &srt.users[j])
}

func (srt *usersSorter) Swap(i, j int) {
	srt.users[i], srt.users[j] = srt.users[j], srt.users[i]
}

func (by By) Sort(users []UserDbRecord, reversed bool) {
	us := &usersSorter{users, by, reversed}
	sort.Sort(us)
}

func OrderBy(users []UserDbRecord, field string, order int) {
	id := func(a, b *UserDbRecord) bool {
		return a.Id < b.Id
	}
	age := func(a, b *UserDbRecord) bool {
		return a.Age < b.Age
	}
	name := func(a, b *UserDbRecord) bool {
		return a.Name < b.Name
	}

	choices := map[string]By {
		"id": id,
		"age": age,
		"name": name,
	}

	if order != 0 {
		By(choices[strings.ToLower(field)]).Sort(users, order == -1)
	}
}

func (db *UsersCollection) Search(req SearchRequest) []UserDbRecord {
	filtered := Filter(db.GetUserRecords(), req.Query)
	OrderBy(filtered, req.OrderField, req.OrderBy)

	return OffsetLimit(filtered, req.Offset, req.Limit)
}

func UserDbRecord2User(u UserDbRecord) (res User) {
	res.Id = u.Id
	res.Name = u.Name
	res.Age = u.Age
	res.About = u.About
	res.Gender = u.Gender
	return
}

func ConvUserDbRecords(recs []UserDbRecord) (users []User){
	for _, u := range recs {
		users = append(users, UserDbRecord2User(u))
	}
	return
}

func extractSearchRequest(r *http.Request) (req SearchRequest, err error) {
	if req.Limit, err = strconv.Atoi(r.FormValue("limit")); err != nil {
		return
	}
	if req.Offset, err = strconv.Atoi(r.FormValue("offset")); err != nil {
		return
	}

	if req.OrderBy, err = strconv.Atoi(r.FormValue("order_by")); err != nil {
		return
	}

	req.Query = r.FormValue("query")
	req.OrderField = r.FormValue("order_field")
	return
}

func SearchServerBadResult(w http.ResponseWriter, r *http.Request) {
	req, _ := extractSearchRequest(r)
	if req.Query == "special_anauthorized" {
		w.WriteHeader(http.StatusUnauthorized)
	} else if req.Query == "special_error" {
		w.WriteHeader(http.StatusInternalServerError)
	} else if req.Query == "special_bad_request1" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "error": 123`))
	} else if req.Query == "special_bad_request2" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "Error": "ErrorBadOrderField" }`))
	} else if req.Query == "special_bad_request3" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "Error": "Unknown" }`))
	} else if req.Query == "special_badjson" {
		w.Write([]byte(`{ "users": 123`))
	} else if req.Query == "special_timeout" {
		time.Sleep(4*time.Second)
	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	var db UsersCollection
	var err error

	if err = db.Load(); err != nil {
		panic(err)
	}

	var req SearchRequest
	if req, err = extractSearchRequest(r); err != nil {

	}

	matched := db.Search(req)

	var data []byte
	if data, err = json.Marshal(ConvUserDbRecords(matched)); err != nil {
		panic(err)
	}
	w.Write(data)
}

func TestSearchClient_FindUsers(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	cli := &SearchClient{URL: ts.URL}
	found, err := cli.FindUsers(SearchRequest{Query: "Mays", Limit: 5})
	assert.Equal(nil, err)
	assert.Equal(1, len(found.Users))
	assert.Equal("Jennings Mays", found.Users[0].Name)
	assert.Equal(false, found.NextPage)
}

func TestSearchClient_FindUsersLimits(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	cli := &SearchClient{URL: ts.URL}
	_, err := cli.FindUsers(SearchRequest{Query: "Mays", Limit: -1})
	assert.Equal(fmt.Errorf("limit must be > 0"), err)
	_, err = cli.FindUsers(SearchRequest{Query: "Mays", Limit: 30})
	assert.Equal(nil, err)
	_, err = cli.FindUsers(SearchRequest{Query: "Mays", Offset: -1})
	assert.Equal(fmt.Errorf("offset must be > 0"), err)
}

// there is no assert package on server, so emulate it
type TestAssert struct {
	t *testing.T
}

func assertNew(t *testing.T) TestAssert {
	return TestAssert{t}
}

func (as *TestAssert) Equal(a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		as.t.Errorf("Comparison failed: %q != %q", a, b)
	}
}

func (as *TestAssert) NotEqual(a, b interface{}) {
	if reflect.DeepEqual(a, b) {
		as.t.Errorf("Comparison failed: %q == %q", a, b)
	}
}

func TestSearchClient_FindUsersNextPage(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	cli := &SearchClient{URL: ts.URL}

	_, err := cli.FindUsers(SearchRequest{Query: "commodo", Limit: 16})
	assert.Equal(nil, err)
}

func TestSearchClient_BadStatuses(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewServer(http.HandlerFunc(SearchServerBadResult))
	cli := &SearchClient{URL: ts.URL}

	_, err := cli.FindUsers(SearchRequest{Query: "special_anauthorized"})
	assert.Equal(fmt.Errorf("Bad AccessToken"), err)

	_, err = cli.FindUsers(SearchRequest{Query: "special_error"})
	assert.Equal(fmt.Errorf("SearchServer fatal error"), err)

	_, err = cli.FindUsers(SearchRequest{Query: "special_bad_request1"})

	_, err = cli.FindUsers(SearchRequest{Query: "special_bad_request2", OrderField: "Game"})
	assert.Equal(fmt.Errorf("OrderFeld %s invalid", "Game"), err)

	_, err = cli.FindUsers(SearchRequest{Query: "special_bad_request3"})
	assert.Equal(fmt.Errorf("unknown bad request error: Unknown"), err)
}

func TestSearchClient_BadBody(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewServer(http.HandlerFunc(SearchServerBadResult))
	cli := &SearchClient{URL: ts.URL}

	_, err := cli.FindUsers(SearchRequest{Query: "special_badjson"})
	assert.NotEqual(nil, err)
}

func TestSearchClient_Disconnected(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewUnstartedServer(http.HandlerFunc(SearchServer))
	cli := &SearchClient{URL: ts.URL}

	_, err := cli.FindUsers(SearchRequest{Query: "commodo"})
	assert.NotEqual(nil, err)
}

func TestSearchClient_Timedout(t *testing.T) {
	assert := assertNew(t)

	ts := httptest.NewServer(http.HandlerFunc(SearchServerBadResult))
	cli := &SearchClient{URL: ts.URL}

	_, err := cli.FindUsers(SearchRequest{Query: "special_timeout"})
	assert.NotEqual(nil, err)
	if !strings.HasPrefix(err.Error(), "timeout") {
		t.Error("should be timed out")
	}
}

func TestLoadUsers(t *testing.T) {
	var db UsersCollection

	if err := db.Load(); err != nil {
		t.Error(err)
	}

	assert := assertNew(t)

	users := db.GetUserRecords()
	assert.Equal(35, len(users))

	assert.Equal("Boyd Wolf", users[0].Name)
	assert.Equal(22, users[0].Age)

	assert.Equal("Kane Sharp", users[len(users)-1].Name)
	assert.Equal(34, users[len(users)-1].Age)
}

func TestFilter(t *testing.T) {
	var db UsersCollection

	if err := db.Load(); err != nil {
		t.Error(err)
	}

	assert := assertNew(t)
	assert.Equal(15, len(db.Search(SearchRequest{Limit: 100, Query: "consectetur"})))
	
	matched := db.Search(SearchRequest{Limit: 2, Query: "Mays"})
	assert.Equal(1, len(matched))
	assert.Equal("Jennings Mays", matched[0].Name)
	assert.Equal(6, matched[0].Id)

	assert.Equal(0, len(db.Search(SearchRequest{Limit: 100, Query: "blahblah"})))
}

func TestFilterLimiting(t *testing.T) {
	var db UsersCollection

	if err := db.Load(); err != nil {
		t.Error(err)
	}

	assert := assertNew(t)
	
	const query = "commodo"

	assert.Equal(17, len(db.Search(SearchRequest{Limit: 100, Query: query})))
	assert.Equal(10, len(db.Search(SearchRequest{Limit: 10, Query: query})))
	assert.Equal(7, len(db.Search(SearchRequest{Offset: 10, Limit: 10, Query: query})))
	assert.Equal(1, len(db.Search(SearchRequest{Offset: 16, Limit: 10, Query: query})))
	assert.Equal(0, len(db.Search(SearchRequest{Offset: 20, Limit: 10, Query: query})))

	// nothing bad with no results
	assert.Equal(0, len(db.Search(SearchRequest{Offset: 50, Limit: 10, Query: "notfound"})))
}

func TestOrdering(t *testing.T) {
	var db UsersCollection

	if err := db.Load(); err != nil {
		t.Error(err)
	}

	assert := assertNew(t)

	{
		res := db.Search(SearchRequest{Limit: 10, Query: "son", OrderField: "Name", OrderBy: 1})
		var names []string
		for _, u := range res {
			names = append(names, u.Name)
		}

		assert.Equal([]string{"Allison Valdez", "Dickson Silva", "Gonzalez Anderson",
			"Henderson Maxwell", "Nicholson Newman", "Whitley Davidson"}, names)
	}

	{
		res := db.Search(SearchRequest{Limit: 10, Query: "son", OrderField: "Name", OrderBy: -1})
		var names []string
		for _, u := range res {
			names = append(names, u.Name)
		}

		assert.Equal([]string{"Whitley Davidson", "Nicholson Newman", "Henderson Maxwell",
		                      "Gonzalez Anderson", "Dickson Silva", "Allison Valdez"}, names)
	}

	{
		res := db.Search(SearchRequest{Limit: 10, Query: "son", OrderField: "Id", OrderBy: -1})
		var ids []int
		for _, u := range res {
			ids = append(ids, u.Id)
		}

		assert.Equal([]int{30, 24, 15, 14, 13, 10}, ids)
	}
}