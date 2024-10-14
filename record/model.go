package record

import (
	"context"

	"github.com/google/uuid"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/sync/_bool"
	"github.com/kyaxcorp/go-helper/sync/_map_string_interface"
	"gorm.io/gorm"
)

const (
	inputDataMapInterface = 1
	inputDataStruct       = 2
)

type Record struct {
	// This is the model which the user has set, after that is the model from which we read back the data!
	// it should be with pointer address
	//NonPtrObj interface{}
	// ================== INPUT DATA ===================\\
	// Data -> can be a map[string]interface{} or a real any Struct from where we will copy the data!
	Data interface{}
	// this case is used when the data is a map[string]interface{}
	dataMap     map[string]interface{}
	dataMapJson string
	// this case is used when the data is structure! but it should be saved as plain structure, because we don't
	// want to modify the input data that came in! So it should be created a copy!
	dataStr interface{}
	// this var is the modelStruct with copied dataStr (all the data from dataStr is copied to dataCopied)
	dataCopied interface{}

	// This is the final data for saving...
	inputData map[string]interface{}

	dataStrHelper *_struct.Helper

	inputDataType int

	// What fields are in input
	// TODO: later on we can replace the string with other type!
	inputFieldNames map[string]string
	// ================== INPUT DATA ===================\\

	//

	// ================== MODEL STRUCTURE ===================\\
	// this is the structure of the set NonPtrObj, it's a pointer!
	ModelStruct interface{}
	// this is not a pointer! it's a plain one!
	modelStruct interface{}
	// ================== MODEL STRUCTURE ===================\\

	//

	// ================== DATABASE DATA ==================== \\
	// this is the same model/data but loaded with existing data from the DB if exists...
	// it will be loaded only if an ID is present (it should be checked on todo: all primary keys) but anyway...
	dbData            interface{}
	loadDataForUpdate bool
	// If it was once loaded
	isDbDataLoaded *_bool.Bool
	// what's the status of the last load
	lastLoadDataStatus *_bool.Bool

	isRecordExists *_bool.Bool
	// ================== DATABASE DATA ==================== \\

	//

	// ================== SAVED DB DATA ==================== \\
	//saveData map[string]interface{}
	// the saveData variable can have or can be supplied with other additional information!
	//saveData map[string]interface{}
	saveData   interface{}
	omitFields map[string]*bool
	// ================== SAVED DB DATA ==================== \\

	//

	// Field Name -> DB column Name
	modelFieldNamings map[string]string

	// AutoLoad -> should auto load the data if it's an existing record
	// TODO: posibil ca acesta nu trebuie sa lucreze in module in care noi cerem datele de la DB apoi facem merge, si apoi
	// actualizam..., ar trebui invers sa facem... sa actualizam doar diferenta in DB, apoi sa cerem record-ul, si apoi
	// sa vedem ce schimbari au fost facute... insa intrebarea consta in aceea -> DAR CARE A FOST OLD RECORD?!
	// Db-ul nu ne raspunde cu astfel de informatie..., de aceea totusi este necesar sa incarcam mai intii informatia
	// aici in memorie si apoi sa actualizam DB doar cu diferenta! diferenta e Input Data!
	// noi am putea sa folosim acelasi input data ca saveData si sa indicam denumirea tabelului manual....
	// sau GORM automat face diferenta?!!
	AutoLoad bool

	Ctx context.Context
	DB  *gorm.DB

	// TODO: add here logger...

	// You can specify: save,create,insert,update
	SaveMode    string
	saveModeSet *_bool.Bool
	//
	dbSet *_bool.Bool

	// it disables default scope for the Loader...
	disableDefaultScope *_bool.Bool

	// these are the record model primary keys
	primaryKeys        map[string]primaryKey
	primaryKeysOrdered []primaryKey

	primaryKeysQuery       string
	primaryKeysItemsValues []interface{}

	cachedIDSet *_bool.Bool
	cachedID    uuid.UUID

	// this is the last error which has been captured
	lastError   error
	lastDBError error

	// this is the userID who made the action!
	userID interface{}
	// this is the deviceID from where the action was made...
	deviceID interface{}

	onBeforeSave *_map_string_interface.MapStringInterface
}
