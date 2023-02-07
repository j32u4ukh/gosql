@ECHO ON
:: 指令格式為: PATH_TO_PROTOCAL_EXECUTE_FILE PATH_TO_PROTOBUF_FILE --csharp_out=PATH_TO_OUTPUT_FOLDER
:: PATH_TO_PROTOCAL_EXECUTE_FILE: 若已將 protoc 編譯程式的路徑加入環境變數，則可直接使用 protoc
protoc ./pb/*.proto --go_out=:./pbgo 
PAUSE