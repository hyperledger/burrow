package sqldb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/burrow/txs"

	"encoding/json"

	"github.com/hyperledger/burrow/vent/types"
)

// findTable checks if a table exists in the default schema
func (db *SQLDB) findTable(tableName string) (bool, error) {

	found := 0
	safeTable := safe(tableName)
	query := clean(db.DBAdapter.FindTableQuery())

	db.Log.Info("msg", "FIND TABLE", "query", query, "value", safeTable)
	if err := db.DB.QueryRow(query, tableName).Scan(&found); err != nil {
		db.Log.Info("msg", "Error finding table", "err", err)
		return false, err
	}

	if found == 0 {
		db.Log.Warn("msg", "Table not found", "value", safeTable)
		return false, nil
	}

	return true, nil
}

// getSysTablesDefinition returns log, chain info & dictionary structures
func (db *SQLDB) getSysTablesDefinition() types.EventTables {

	tables := make(types.EventTables)
	dicCol := make(map[string]types.SQLTableColumn)
	logCol := make(map[string]types.SQLTableColumn)
	chainCol := make(map[string]types.SQLTableColumn)
	errorsCol := make(map[string]types.SQLTableColumn)

	// log table
	logCol[types.SQLColumnLabelId] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelId,
		Type:    types.SQLColumnTypeSerial,
		Primary: true,
		Order:   1,
	}

	logCol[types.SQLColumnLabelTimeStamp] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelTimeStamp,
		Type:    types.SQLColumnTypeTimeStamp,
		Primary: false,
		Order:   2,
	}

	logCol[types.SQLColumnLabelTableName] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelTableName,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: false,
		Order:   3,
	}

	logCol[types.SQLColumnLabelEventName] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelEventName,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: false,
		Order:   4,
	}

	logCol[types.SQLColumnLabelEventFilter] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelEventFilter,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: false,
		Order:   5,
	}

	logCol[types.SQLColumnLabelHeight] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelHeight,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: false,
		Order:   6,
	}

	logCol[types.SQLColumnLabelTxHash] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelTxHash,
		Type:    types.SQLColumnTypeVarchar,
		Length:  txs.HashLengthHex,
		Primary: false,
		Order:   7,
	}

	logCol[types.SQLColumnLabelAction] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelAction,
		Type:    types.SQLColumnTypeVarchar,
		Length:  20,
		Primary: false,
		Order:   8,
	}

	logCol[types.SQLColumnLabelDataRow] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelDataRow,
		Type:    types.SQLColumnTypeJSON,
		Length:  0,
		Primary: false,
		Order:   9,
	}

	logCol[types.SQLColumnLabelSqlStmt] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelSqlStmt,
		Type:    types.SQLColumnTypeText,
		Length:  0,
		Primary: false,
		Order:   10,
	}

	logCol[types.SQLColumnLabelSqlValues] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelSqlValues,
		Type:    types.SQLColumnTypeText,
		Length:  0,
		Primary: false,
		Order:   11,
	}

	// dictionary table
	dicCol[types.SQLColumnLabelTableName] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelTableName,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: true,
		Order:   1,
	}

	dicCol[types.SQLColumnLabelColumnName] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelColumnName,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: true,
		Order:   2,
	}

	dicCol[types.SQLColumnLabelColumnType] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelColumnType,
		Type:    types.SQLColumnTypeInt,
		Length:  0,
		Primary: false,
		Order:   3,
	}

	dicCol[types.SQLColumnLabelColumnLength] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelColumnLength,
		Type:    types.SQLColumnTypeInt,
		Length:  0,
		Primary: false,
		Order:   4,
	}

	dicCol[types.SQLColumnLabelPrimaryKey] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelPrimaryKey,
		Type:    types.SQLColumnTypeInt,
		Length:  0,
		Primary: false,
		Order:   5,
	}

	dicCol[types.SQLColumnLabelColumnOrder] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelColumnOrder,
		Type:    types.SQLColumnTypeInt,
		Length:  0,
		Primary: false,
		Order:   6,
	}

	// chain info table
	chainCol[types.SQLColumnLabelChainID] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelChainID,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: true,
		Order:   1,
	}

	chainCol[types.SQLColumnLabelBurrowVer] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelBurrowVer,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: false,
		Order:   2,
	}

	// errors info table
	errorsCol[types.SQLColumnLabelErrorsHeight] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelErrorsHeight,
		Type:    types.SQLColumnTypeVarchar,
		Length:  100,
		Primary: false,
		Order:   1,
	}

	errorsCol[types.SQLColumnLabelErrorsTxHash] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelErrorsTxHash,
		Type:    types.SQLColumnTypeVarchar,
		Length:  txs.HashLengthHex,
		Primary: true,
		Order:   2,
	}

	errorsCol[types.SQLColumnLabelErrorsMessage] = types.SQLTableColumn{
		Name:    types.SQLColumnLabelErrorsMessage,
		Type:    types.SQLColumnTypeVarchar,
		Length:  300,
		Primary: false,
		Order:   3,
	}

	// add tables
	//log
	tables[types.SQLLogTableName] = types.SQLTable{
		Name:    types.SQLLogTableName,
		Columns: logCol,
	}

	//dictionary
	tables[types.SQLDictionaryTableName] = types.SQLTable{
		Name:    types.SQLDictionaryTableName,
		Columns: dicCol,
	}

	//chain info
	tables[types.SQLChainInfoTableName] = types.SQLTable{
		Name:    types.SQLChainInfoTableName,
		Columns: chainCol,
	}

	tables[types.SQLErrorsTableName] = types.SQLTable{
		Name:    types.SQLErrorsTableName,
		Columns: errorsCol,
	}

	return tables
}

// getTableDef returns the structure of a given SQL table
func (db *SQLDB) getTableDef(tableName string) (types.SQLTable, error) {

	var table types.SQLTable

	safeTable := safe(tableName)

	found, err := db.findTable(safeTable)
	if err != nil {
		return table, err
	}

	if !found {
		db.Log.Info("msg", "Error table not found", "value", safeTable)
		return table, errors.New("Error table not found " + safeTable)
	}

	table.Name = safeTable
	query := clean(db.DBAdapter.TableDefinitionQuery())

	db.Log.Info("msg", "QUERY STRUCTURE", "query", query, "value", safeTable)
	rows, err := db.DB.Query(query, safeTable)
	if err != nil {
		db.Log.Info("msg", "Error querying table structure", "err", err)
		return table, err
	}
	defer rows.Close()

	columns := make(map[string]types.SQLTableColumn)
	i := 0

	for rows.Next() {
		var columnName string
		var columnSQLType types.SQLColumnType
		var columnIsPK int
		var columnLength int
		var column types.SQLTableColumn

		if err = rows.Scan(&columnName, &columnSQLType, &columnLength, &columnIsPK); err != nil {
			db.Log.Info("msg", "Error scanning table structure", "err", err)
			return table, err
		}

		if _, err = db.DBAdapter.TypeMapping(columnSQLType); err != nil {
			return table, err
		}

		column.Name = columnName
		column.Type = columnSQLType
		column.Length = columnLength
		column.Primary = columnIsPK == 1
		column.Order = i

		columns[columnName] = column
		i++
	}

	if err = rows.Err(); err != nil {
		db.Log.Info("msg", "Error during rows iteration", "err", err)
		return table, err
	}

	table.Columns = columns
	return table, nil
}

// alterTable alters the structure of a SQL table & add info to the dictionary
func (db *SQLDB) alterTable(newTable types.SQLTable, eventName string) error {

	db.Log.Info("msg", "Altering table", "value", newTable.Name)

	// prepare log query
	logQuery := clean(db.DBAdapter.InsertLogQuery())

	// current table structure
	safeTable := safe(newTable.Name)
	currentTable, err := db.getTableDef(safeTable)
	if err != nil {
		return err
	}

	sqlValues, _ := db.getJSON(nil)

	// for each column in the new table structure
	for _, newColumn := range newTable.Columns {
		found := false

		// check if exists in the current table structure
		for _, currentColumn := range currentTable.Columns {
			// if column exists
			if currentColumn.Name == newColumn.Name {
				found = true
				break
			}
		}

		if !found {
			safeCol := safe(newColumn.Name)
			query, dictionary := db.DBAdapter.AlterColumnQuery(safeTable, safeCol, newColumn.Type, newColumn.Length, newColumn.Order)

			//alter column
			db.Log.Info("msg", "ALTER TABLE", "query", safe(query))
			_, err = db.DB.Exec(safe(query))

			if err != nil {
				if db.DBAdapter.ErrorEquals(err, types.SQLErrorTypeDuplicatedColumn) {
					db.Log.Warn("msg", "Duplicate column", "value", safeCol)
				} else {
					db.Log.Info("msg", "Error altering table", "err", err)
					return err
				}
			} else {
				//store dictionary
				db.Log.Info("msg", "STORE DICTIONARY", "query", clean(dictionary))
				_, err = db.DB.Exec(dictionary)
				if err != nil {
					db.Log.Info("msg", "Error storing  dictionary", "err", err)
					return err
				}

				//insert log (if action is not database initialization)
				if eventName != string(types.ActionInitialize) {
					// Marshal the table into a JSON string.
					var jsonData []byte
					jsonData, err = db.getJSON(newColumn)
					if err != nil {
						db.Log.Info("msg", "error marshaling column", "err", err, "value", fmt.Sprintf("%v", newColumn))
						return err
					}
					//insert log
					_, err = db.DB.Exec(logQuery, newTable.Name, eventName, newTable.Filter, nil, nil, types.ActionAlterTable, jsonData, query, sqlValues)
					if err != nil {
						db.Log.Info("msg", "Error inserting log", "err", err)
						return err
					}
				}
			}
		}
	}
	return nil
}

// getSelectQuery builds a select query for a specific SQL table and a given block
func (db *SQLDB) getSelectQuery(table types.SQLTable, height string) (string, error) {

	fields := ""

	for _, tableColumn := range table.Columns {
		if fields != "" {
			fields += ", "
		}
		fields += db.DBAdapter.SecureColumnName(tableColumn.Name)
	}

	if fields == "" {
		return "", errors.New("error table does not contain any fields")
	}

	query := clean(db.DBAdapter.SelectRowQuery(table.Name, fields, height))
	return query, nil
}

// createTable creates a new table
func (db *SQLDB) createTable(table types.SQLTable, eventName string) error {

	db.Log.Info("msg", "Creating Table", "value", table.Name)

	// prepare log query
	logQuery := clean(db.DBAdapter.InsertLogQuery())

	// sort columns
	columns := len(table.Columns)
	sortedColumns := make([]types.SQLTableColumn, columns)
	for _, tableColumn := range table.Columns {
		if tableColumn.Order <= 0 {
			db.Log.Info("msg", "column_order <=0")
			return fmt.Errorf("table definition error,%s has column_order <=0 (minimum value = 1)", tableColumn.Name)
		} else if tableColumn.Order-1 > columns {
			db.Log.Info("msg", "column_order > total_columns")
			return fmt.Errorf("table definition error, %s has column_order > total_columns", tableColumn.Name)
		} else if sortedColumns[tableColumn.Order-1].Order != 0 {
			db.Log.Info("msg", "duplicated column_oder")
			return fmt.Errorf("table definition error, %s and %s have duplicated column_order", sortedColumns[tableColumn.Order-1].Name, tableColumn.Name)
		} else {
			sortedColumns[tableColumn.Order-1] = tableColumn
		}
	}

	//get create table query
	safeTable := safe(table.Name)
	query, dictionary := db.DBAdapter.CreateTableQuery(safeTable, sortedColumns)
	if query == "" {
		db.Log.Info("msg", "empty CREATE TABLE query")
		return errors.New("empty CREATE TABLE query")
	}

	// create table
	db.Log.Info("msg", "CREATE TABLE", "query", clean(query))
	_, err := db.DB.Exec(query)
	if err != nil {
		if db.DBAdapter.ErrorEquals(err, types.SQLErrorTypeDuplicatedTable) {
			db.Log.Warn("msg", "Duplicate table", "value", safeTable)
			return nil

		} else if db.DBAdapter.ErrorEquals(err, types.SQLErrorTypeInvalidType) {
			db.Log.Info("msg", "Error creating table, invalid datatype", "err", err)
			return err

		}
		db.Log.Info("msg", "Error creating table", "err", err)
		return err
	}

	//store dictionary
	db.Log.Info("msg", "STORE DICTIONARY", "query", clean(dictionary))
	_, err = db.DB.Exec(dictionary)
	if err != nil {
		db.Log.Info("msg", "Error storing  dictionary", "err", err)
		return err
	}

	//insert log (if action is not database initialization)
	if eventName != string(types.ActionInitialize) {

		// Marshal the table into a JSON string.
		var jsonData []byte
		jsonData, err = db.getJSON(table)
		if err != nil {
			db.Log.Info("msg", "error marshaling table", "err", err, "value", fmt.Sprintf("%v", table))
			return err
		}
		sqlValues, _ := db.getJSON(nil)

		//insert log
		_, err = db.DB.Exec(logQuery, table.Name, eventName, table.Filter, nil, nil, types.ActionCreateTable, jsonData, query, sqlValues)
		if err != nil {
			db.Log.Info("msg", "Error inserting log", "err", err)
			return err
		}
	}
	return nil
}

// getBlockTables return all SQL tables that have been involved
// in a given batch transaction for a specific block
func (db *SQLDB) getBlockTables(block string) (types.EventTables, error) {

	tables := make(types.EventTables)

	query := clean(db.DBAdapter.SelectLogQuery())
	db.Log.Info("msg", "QUERY LOG", "query", query, "value", block)

	rows, err := db.DB.Query(query, block)
	if err != nil {
		db.Log.Info("msg", "Error querying log", "err", err)
		return tables, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventName, tableName string
		var table types.SQLTable

		err = rows.Scan(&tableName, &eventName)
		if err != nil {
			db.Log.Info("msg", "Error scanning table structure", "err", err)
			return tables, err
		}

		err = rows.Err()
		if err != nil {
			db.Log.Info("msg", "Error scanning table structure", "err", err)
			return tables, err
		}

		table, err = db.getTableDef(tableName)
		if err != nil {
			return tables, err
		}

		tables[eventName] = table
	}

	return tables, nil
}

// clean queries from tabs, spaces  and returns
func clean(parameter string) string {
	replacer := strings.NewReplacer("\n", " ", "\t", "")
	return replacer.Replace(parameter)
}

// safe sanitizes a parameter
func safe(parameter string) string {
	replacer := strings.NewReplacer(";", "", ",", "")
	return replacer.Replace(parameter)
}

//getJSON returns marshaled json from JSON single column
func (db *SQLDB) getJSON(JSON interface{}) ([]byte, error) {
	if JSON != nil {
		return json.Marshal(JSON)
	}
	return json.Marshal("")
}

//getJSONFromValues returns marshaled json from query values
func (db *SQLDB) getJSONFromValues(values []interface{}) ([]byte, error) {
	if values != nil {
		return json.Marshal(values)
	}
	return json.Marshal("")
}

//getValuesFromJSON returns query values from unmarshaled JSON column
func (db *SQLDB) getValuesFromJSON(JSON string) ([]interface{}, error) {
	pointers := make([]interface{}, 0)
	bytes := []byte(JSON)
	err := json.Unmarshal(bytes, &pointers)
	return pointers, err
}
