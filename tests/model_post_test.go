package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/victorsteven/fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllDrugs(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatalf("Error refreshing user and drugs table %v\n", err)
	}
	_, _, err = seedDataTable()
	if err != nil {
		log.Fatalf("Error seeding user and drug  table %v\n", err)
	}
	drugs, err := postInstance.FindAllDrug(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*drugs), 2)
}

func TestSaveDrug(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	newPost := models.Drug{
		ID:       1,
		Name:     "Aspirina",
		Approved: true,
		Min_dose: 1,
		Max_dose: 2,
	}
	savedDrug, err := newPost.SaveDrug(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newPost.ID, savedDrug.ID)
	assert.Equal(t, newPost.Name, savedDrug.Name)
	assert.Equal(t, newPost.Approved, savedDrug.Approved)
	assert.Equal(t, newPost.Min_dose, savedDrug.Min_dose)
	assert.Equal(t, newPost.Max_dose, savedDrug.Max_dose)
}

func TestGetDrugByID(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	drug, err := seedOneDataTable()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundDrug, err := postInstance.FindDrugByID(server.DB, drug.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundDrug.ID, drug.ID)
	assert.Equal(t, foundDrug.Name, drug.Name)
	assert.Equal(t, foundDrug.Approved, drug.Approved)
	assert.Equal(t, foundDrug.Min_dose, drug.Min_dose)
	assert.Equal(t, foundDrug.Max_dose, drug.Max_dose)
}

func TestUpdateADrug(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneDataTable()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	drugUpdate := models.Drug{
		ID:       1,
		Name:     "Aspirina",
		Approved: true,
		Min_dose: 1,
		Max_dose: 2,
	}
	updatedPost, err := drugUpdate.UpdateADrug(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedPost.ID, drugUpdate.ID)
	assert.Equal(t, updatedPost.Name, drugUpdate.Name)
	assert.Equal(t, updatedPost.Approved, drugUpdate.Approved)
	assert.Equal(t, updatedPost.Min_dose, drugUpdate.Min_dose)
	assert.Equal(t, updatedPost.Max_dose, drugUpdate.Max_dose)
}

func TestDeleteADrug(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	drug, err := seedOneDataTable()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := postInstance.DeleteADrug(server.DB, drug.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
