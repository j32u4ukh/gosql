# Stmt

**目標：**

* 盡可能不使用反射或是第三方套件，純粹的使用字串拼接，產生 SQL 所需語法指令。
* 所有指令，就算操作上會比較麻煩，stmt 也應該都能生成出來。
* stmt 對於傳入參數與配置採相信為原則，因此使用者須確保傳入參數的正確性。
* 追求易用性的使用者，可以使用對 stmt 做過進一步封裝的套件。由於該套件可能會使用反射等操作，導致效能略為下降，但使用上更為方便。

# ColumnParam 附加定義

1. default: AI(Auto Increment), NULL, NIL(不設置預設值), 其他(current_timestamp())
2. size: 定義 DB 欄位變數大小([SQL 類型]: [變數大小]) 數值型用預設值，字串型沒有值，僅允許 VARCHAR 設定
> * TINYINT(4):     -128 ~ 127 / 0 ~ 255
> * SMALINT(6):     -3,2768 ~ 3,2767 / 0 ~ 6,5535
> * MEDIUMINT(9):   -838,8608 ~ 838,8607 / 0 ~ 1677,7215
> * INT(11):        -21,4748,3648 ~ 21,4748,3647 / 0 ~ 42,9496,7295
> * BIGINT(20):     -922,3372,0368,5477,5808 ~ 922,3372,0368,5477,5807 / 0 ~ 1844,6744,0737,0955,1615 
> * TINYTEXT(-):    255
> * VARCHAR(3000):  0 ~ 6,5535
> * TEXT(-):        6,5535
> * MEDIUMTEXT(-):  1677,7215
> * LONGTEXT(-):    42,9496,7295
3. type: 定義 DB 欄位變數類型([SQL 類型]: [Proto 類型])
> * TINYINT:    int32/bool
> * SMALINT:    int32
> * MEDIUMINT:  int32
> * INT:        int32
> * BIGINT:     int64
> * TINYTEXT:   string
> * VARCHAR:    string/map/Message(TIMESTAMP 以外)
> * TEXT:       string
> * MEDIUMTEXT: string
> * LONGTEXT:   string
> * TIMESTAMP(-):  TIMESTAMP(Message)
4. can_null: 是否可以為空值 (true: 可以 / false: 不可以) **不給設
5. primary_key: 是否為主鍵，數值為演算法(填入 default 則使用預設值)
6. comment: 註解內容
7. unsigned: 沒有負數？ ("unsigned": "true" / "false" )
8. update: ON UPDATE 時執行的函數

* 字串中需要空格可以使用 \t

## 調整 DB 欄位順序，應與 Proto 檔一致
---

建議參數修改順序如下：

1. Comment
2. CanNull
3. Type
4. Size
5. Unsigned
6. Default
7. PrimaryKey
8. Update