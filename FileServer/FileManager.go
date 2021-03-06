package FileServer

import (
	"duov6.com/FileServer/messaging"
	"duov6.com/common"
	"duov6.com/objectstore/client"
	goclient "duov6.com/objectstore/client"
	"duov6.com/objectstore/repositories"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/toqueteos/webbrowser"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type FileManager struct {
}

type FileData struct {
	Id       string
	FileName string
	Body     string
}

var uploadFileName string

func (f *FileManager) Store(request *messaging.FileRequest) messaging.FileResponse { // store disk on database

	fileResponse := messaging.FileResponse{}

	if len(request.Body) == 0 {

		//WHEN REQUEST COMES FROM A REST INTERFACE
		file, header, err := request.WebRequest.FormFile("file")

		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
			return fileResponse
		}

		if header == nil {
			fmt.Println("No Header Found!")
			uploadFileName = request.Parameters["id"]
		} else {
			uploadFileName = header.Filename
		}

		out, err := os.Create(uploadFileName)
		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
			return fileResponse
		}

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
			return fileResponse
		}

		file2, err2 := ioutil.ReadFile(uploadFileName)

		if err2 != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
			return fileResponse
		}

		convertedBody := string(file2[:])
		base64Body := common.EncodeToBase64(convertedBody)

		//Create a instance of file struct
		obj := FileData{}
		obj.Id = request.Parameters["id"]
		obj.FileName = uploadFileName
		obj.Body = base64Body

		var extraMap map[string]interface{}
		extraMap = make(map[string]interface{})
		extraMap["File"] = "excelFile"

		fmt.Println("Namespace : " + request.Parameters["namespace"])
		fmt.Println("Class : " + request.Parameters["class"])

		uploadContext := strings.ToLower(request.Parameters["fileContent"])

		isRawFile := false
		isIndividualData := false

		if uploadContext == "" || uploadContext == "both" || uploadContext == "raw" {
			isRawFile = true
		}
		if uploadContext == "" || uploadContext == "both" || uploadContext == "data" {
			isIndividualData = true
		}

		if isIndividualData {
			fmt.Println("Saving INDIVIDUAL DATA inside file.......... ")
			if checkIfFile(uploadFileName) == "xlsx" {
				isRawFile = false
				status := SaveExcelEntries(uploadFileName, request)
				if status == true {
					fmt.Println("Individual Records Saved Successfully!")
				} else {
					fmt.Println("Saving Individual Records Failed!")
				}
			}
		}

		var returnParams repositories.RepositoryResponse
		if isRawFile {
			fmt.Println("Saving the RAW file.......... ")
			returnParams = client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], request.Parameters["class"], extraMap).StoreObject().WithKeyField("Id").AndStoreOne(obj).FileOk()
			if len(returnParams.Data) > 0 {
				fmt.Fprintf(request.WebResponse, returnParams.Data[0]["ID"].(string))
			} else {
				fmt.Fprintf(request.WebResponse, "FAILED!")
			}
		} else {
			fmt.Fprintf(request.WebResponse, uploadFileName)
		}

		//close the files
		err = out.Close()

		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
			return fileResponse
		}

		err = file.Close()

		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
			return fileResponse
		}

		//remove the temporary stored file from the disk
		err2 = os.Remove(uploadFileName)

		if err2 != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err2.Error()
			return fileResponse
		}

		if err == nil && err2 == nil {
			fileResponse.IsSuccess = true
			fileResponse.Message = "Storing file successfully completed"
		} else {
			fileResponse.IsSuccess = false
			fileResponse.Message = "Storing file was unsuccessfull!" + "\n" + err.Error() + "\n" + err2.Error()
		}

	} else {

		//WHEN REQUEST COMES FROM A NON REST INTERFACE
		convertedBody := string(request.Body[:])
		base64Body := common.EncodeToBase64(convertedBody)

		//store file in the DB as a single file
		obj := FileData{}
		obj.Id = request.Parameters["id"]
		obj.FileName = request.FileName
		obj.Body = base64Body

		response := client.Go(request.Parameters["securityToken"], request.Parameters["namespace"], request.Parameters["class"]).StoreObject().WithKeyField("Id").AndStoreOne(obj).FileOk()
		fileResponse.IsSuccess = response.IsSuccess
		fileResponse.Message = response.Message
	}

	return fileResponse
}

func (f *FileManager) Remove(request *messaging.FileRequest) messaging.FileResponse { // remove file from disk and database
	fileResponse := messaging.FileResponse{}

	//Delete from Physical location
	var saveServerPath string = request.RootSavePath
	file, err := ioutil.ReadFile(saveServerPath + request.FilePath + request.FileName)

	if err != nil {
		if len(file) > 0 {
			err = os.Remove(saveServerPath + request.FilePath + request.FileName)
		}
	} else {
		fmt.Println("Physical file not available to Delete!")
	}

	//Delete from ObjectStore
	obj := FileData{}
	obj.Id = request.Parameters["id"]
	obj.FileName = request.FileName

	err = goclient.Go(request.Parameters["securityToken"], request.Parameters["namespace"], request.Parameters["class"]).StoreObjectWithOperation("delete").WithKeyField("Id").AndStoreOne(obj).Ok()

	if err == nil {
		fileResponse.IsSuccess = true
		fileResponse.Message = "Deletion of file successfully completed"
	} else {
		fileResponse.IsSuccess = true
		fileResponse.Message = err.Error()
	}

	return fileResponse
}

func (f *FileManager) Download(request *messaging.FileRequest) messaging.FileResponse { // save the file to ftp and download via ftp on browser
	fileResponse := messaging.FileResponse{}

	if len(request.Body) != 0 {
		var saveServerPath string = request.RootSavePath
		var accessServerPath string = request.RootGetPath

		file := FileData{}
		json.Unmarshal(request.Body, &file)

		temp := common.DecodeFromBase64(file.Body)
		ioutil.WriteFile((saveServerPath + request.FilePath + file.FileName), []byte(temp), 0666)
		err := webbrowser.Open(accessServerPath + request.FilePath + file.FileName)
		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
		} else {
			fileResponse.IsSuccess = true
			fileResponse.Message = "Downloading file successfully completed"
		}
	} else {
		fileResponse.IsSuccess = false
		fileResponse.Message = "No Request Body Found!"
	}

	return fileResponse
}

func SaveExcelEntries(excelFileName string, request *messaging.FileRequest) bool {
	fmt.Println("Inserting Records to Database....")
	rowcount := 0
	colunmcount := 0
	var exceldata []map[string]interface{}
	var colunName []string

	//file read
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
				}
				break
			}
			exceldata = make(([]map[string]interface{}), rowcount)
			if error == nil {
				for _, sheet := range xlFile.Sheets {
					for rownumber, row := range sheet.Rows {
						currentRow := make(map[string]interface{})
						if rownumber != 0 {
							exceldata[rownumber-1] = currentRow
							for cellnumber, cell := range row.Cells {
								if cellnumber == 0 {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 0 {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 2 {
									dd, _ := cell.Float()
									exceldata[rownumber-1][colunName[cellnumber]] = float64(dd)
								} else if cell.Type() == 3 {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.Bool()
								} else {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								}
							}
						}
					}
				}
			}

			Id := colunName[0]
			var extraMap map[string]interface{}
			extraMap = make(map[string]interface{})
			extraMap["File"] = "exceldata"
			fmt.Println("Namespace : " + request.Parameters["namespace"])
			fmt.Println("Keyfield : " + Id)
			fmt.Println("filename : " + getExcelFileName(excelFileName))
			client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], getExcelFileName(excelFileName), extraMap).StoreObject().WithKeyField(Id).AndStoreMapInterface(exceldata).Ok()
			return true
		}

	}
	return false
}

func checkIfFile(params string) (fileType string) {
	var tempArray []string
	tempArray = strings.Split(params, ".")
	if len(tempArray) > 1 {
		fileType = tempArray[len(tempArray)-1]
	} else {
		fileType = "NAF"
	}
	return
}

func getExcelFileName(path string) (fileName string) {
	subsets := strings.Split(path, "\\")
	subfilenames := strings.Split(subsets[len(subsets)-1], ".")
	fileName = subfilenames[0]
	return
}
