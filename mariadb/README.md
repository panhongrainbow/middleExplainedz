# MariaDB SP Explanation

## Variable Lifecycle

### Initialization settings

```bash
$ mysql -u root -P 3306 -p # Password 12345
$ CREATE USER 'panhong'@'localhost' IDENTIFIED BY '12345';
$ GRANT ALL PRIVILEGES ON *.* TO 'panhong'@'localhost' WITH GRANT OPTION;
$ FLUSH PRIVILEGES;

$ mysql -u panhong -h localhost -P 3306 -p # Password 12345
$ CREATE DATABASE panhong;
$ USE panhong;
```

### Check SP status

```bash
$ SHOW PROCEDURE STATUS WHERE Db = 'panhong' AND Name NOT LIKE 'mysql%';

$ DROP PROCEDURE var_Lifecycle;

$ CALL var_Lifecycle();
```

### Generate a test SP

```sql
-- Set the delimiter // 
DELIMITER //

-- Generate SP, named var_Lifecycle()
CREATE PROCEDURE var_Lifecycle()
-- SP starts
BEGIN
    -- Declaring variables of different types
    DECLARE count INT DEFAULT 0;
    SET count = 10;
    SET @exist = 11;

-- SP ends
END//

-- Restore delimiter ;
DELIMITER ;
```

### Call the test SP

```sql
# >>>>> >>>>> >>>>> >>>>> >>>>> in MariaDB Section1

$ mysql -u panhong -h localhost -P 3306 -p # Password 12345

USE panhong;

CALL var_Lifecycle();
# Query OK, 0 rows affected (0.001 sec)

SELECT @exist;
#+--------+
#| @exist |
#+--------+
#|     11 |
#+--------+
# 1 row in set (0.000 sec)

SELECT count;
# ERROR 1054 (42S22): Unknown column 'count' in 'field list'



# >>>>> >>>>> >>>>> >>>>> >>>>> in MariaDB Section2

$ mysql -u panhong -h localhost -P 3306 -p # Password 12345

USE panhong;

SELECT @exist;
#+--------+
#| @exist |
#+--------+
#| NULL   |
#+--------+
#1 row in set (0.000 sec)



# >>>>> >>>>> >>>>> >>>>> >>>>> in MariaDB Section1
exit

mysql -u panhong -h localhost -P 3306 -p # Password 12345

USE panhong

SELECT @exist;
#+--------+
#| @exist |
#+--------+
#| NULL   |
#+--------+
#1 row in set (0.000 sec)
```

As can be seen from the above, the @exist variable `only exists in the MariaDB session that was logged in`.

After logging out, the @exist variable will `be released`. 

(只存在于登入的MariaDB 的 Session 里，登出后被适放)

## Generate an array using JSON

### Initialization settings

```bash
$ mysql -u root -P 3306 -p # Password 12345
$ CREATE USER 'panhong'@'localhost' IDENTIFIED BY '12345';
$ GRANT ALL PRIVILEGES ON *.* TO 'panhong'@'localhost' WITH GRANT OPTION;
$ FLUSH PRIVILEGES;

$ mysql -u panhong -h localhost -P 3306 -p # Password 12345
$ CREATE DATABASE panhong;
$ USE panhong;
```

### Check SP status

```bash
$ SHOW PROCEDURE STATUS WHERE Db = 'panhong' AND Name NOT LIKE 'mysql%';

$ DROP PROCEDURE array_from_json;

$ CALL array_from_json();
```

### Generate an array on SP

```sql
-- Set the delimiter // 
DELIMITER //

-- Generate SP, named array_from_json()
CREATE PROCEDURE array_from_json()
-- SP starts
BEGIN
    -- Declare month Json variable
    DECLARE months JSON DEFAULT NULL;
    SET months = JSON_ARRAY(
        JSON_OBJECT('Id', 1, 'MonthName', 'January'),
        JSON_OBJECT('Id', 2, 'MonthName', 'February'),
        JSON_OBJECT('Id', 3, 'MonthName', 'March'),
        JSON_OBJECT('Id', 4, 'MonthName', 'April'),
        JSON_OBJECT('Id', 5, 'MonthName', 'May'),
        JSON_OBJECT('Id', 6, 'MonthName', 'June'),
        JSON_OBJECT('Id', 7, 'MonthName', 'July'),
        JSON_OBJECT('Id', 8, 'MonthName', 'August'),
        JSON_OBJECT('Id', 9, 'MonthName', 'September'),
        JSON_OBJECT('Id', 10, 'MonthName', 'October'),
        JSON_OBJECT('Id', 11, 'MonthName', 'November'),
        JSON_OBJECT('Id', 12, 'MonthName', 'December')
    );
    
    -- Return month array
    SELECT JSON_EXTRACT(months, CONCAT('$[', Id - 1, '].MonthName')) AS MonthArray FROM (SELECT 1 AS Id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10 UNION SELECT 11 UNION SELECT 12) AS T;

-- SP ends
END//

-- Restore delimiter ;
DELIMITER ;
```

### Return month array

#### JSON_EXTRACT function

Get part of the JSON content

```bash
SELECT JSON_EXTRACT('{"name": "John", "age": 30}', '$.name');
#+-------------------------------------------------------+
#| JSON_EXTRACT('{"name": "John", "age": 30}', '$.name') |
#+-------------------------------------------------------+
#| "John"                                                |
#+-------------------------------------------------------+
#1 row in set (0.001 sec)

SELECT JSON_EXTRACT('{"name": "John", "age": 30}', '$[0].name');
#+----------------------------------------------------------+
#| JSON_EXTRACT('{"name": "John", "age": 30}', '$[0].name') |
#+----------------------------------------------------------+
#| "John"                                                   |
#+----------------------------------------------------------+
#1 row in set (0.001 sec)

SELECT JSON_EXTRACT('{"name": "John", "age": 30}', '$[1].name');
#+----------------------------------------------------------+
#| JSON_EXTRACT('{"name": "John", "age": 30}', '$[1].name') |
#+----------------------------------------------------------+
#| NULL                                                     |
#+----------------------------------------------------------+
#1 row in set (0.001 sec)
```

#### UNIT function

Combine all the query results, the following will return 1 to 12 arrays

As below, combine the results of SELECT 1 and SELECT 2 ...

```sql
SELECT 1 AS Id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10 UNION SELECT 11 UNION SELECT 12;
#+----+
#| Id |
#+----+
#|  1 |
#|  2 |
#|  3 |
#|  4 |
#|  5 |
#|  6 |
#|  7 |
#|  8 |
#|  9 |
#| 10 |
#| 11 |
#| 12 |
#+----+
#12 rows in set (0.001 sec)
```

#### CONCAT function

Used for string concatenation as follows

````sql
SELECT CONCAT('hello', 'world');
#+--------------------------+
#| CONCAT('hello', 'world') |
#+--------------------------+
#| helloworld               |
#+--------------------------+
#1 row in set (0.000 sec)
````

#### Full SQL sentence

```sql
SELECT JSON_EXTRACT(months, CONCAT('$[', Id - 1, '].MonthName')) AS MonthArray FROM (SELECT 1 AS Id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10 UNION SELECT 11 UNION SELECT 12) AS T;

# This sentence can be executed independently
# (SELECT 1 AS Id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10 UNION SELECT 11 UNION SELECT 12)
# Will generate a temporary table from 1 to 12, named T
# Because this sentence has SELECT 1 AS Id, so after that Id will change from 1 to 12

CONCAT('$[', Id - 1, '].MonthName')
# This sentence cannot be executed independently
# Because Id changes from 1 to 12, the following contents will be produced
# '$[0].MonthName'
# '$[1].MonthName'
# '$[2].MonthName'
# '$[3].MonthName'
# '$[4].MonthName'
# '$[5].MonthName'
# '$[6].MonthName'
# '$[7].MonthName'
# '$[8].MonthName'
# '$[9].MonthName'
# '$[10].MonthName'
# '$[11].MonthName'

SELECT JSON_EXTRACT(months, CONCAT('$[', Id - 1, '].MonthName'))
# This sentence cannot be executed independently
# This sentence will become 12 SQL subclauses
# SELECT JSON_EXTRACT(months, '$[0].MonthName')
# SELECT JSON_EXTRACT(months, '$[1].MonthName')
# SELECT JSON_EXTRACT(months, '$[2].MonthName')
# SELECT JSON_EXTRACT(months, '$[3].MonthName')
# SELECT JSON_EXTRACT(months, '$[4].MonthName')
# SELECT JSON_EXTRACT(months, '$[5].MonthName')
# SELECT JSON_EXTRACT(months, '$[6].MonthName')
# SELECT JSON_EXTRACT(months, '$[7].MonthName')
# SELECT JSON_EXTRACT(months, '$[8].MonthName')
# SELECT JSON_EXTRACT(months, '$[9].MonthName')
# SELECT JSON_EXTRACT(months, '$[10].MonthName')
# SELECT JSON_EXTRACT(months, '$[11].MonthName')
# Just take the 12 month names from the Json array
```

## Generate an array using Case

The content is as follows

````sql
SELECT CASE Id
    WHEN 1 THEN 'January'
    WHEN 2 THEN 'February'
    WHEN 3 THEN 'March'
    WHEN 4 THEN 'April'
    WHEN 5 THEN 'May'
    WHEN 6 THEN 'June'
    WHEN 7 THEN 'July'
    WHEN 8 THEN 'August'
    WHEN 9 THEN 'September'
    WHEN 10 THEN 'October'
    WHEN 11 THEN 'November'
    WHEN 12 THEN 'December'
END AS MonthArray FROM (
    SELECT 1 AS Id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10 UNION SELECT 11 UNION SELECT 12
) AS T;
````







