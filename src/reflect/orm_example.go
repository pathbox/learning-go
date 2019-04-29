package orm

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int64     `json:"id"`         // 自增主键
	Age       int64     `json:"age"`        // 年龄
	FirstName string    `json:"first_name"` // 姓
	LastName  string    `json:"last_name"`  // 名
	Email     string    `json:"email"`      // 邮箱地址
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

func Connect(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(10 * time.Minute)
	return conn, conn.Ping()
}

//Query will build a sql
type Query struct {
	db      *sql.DB
	table   string
}

//Table bind db and table
func Table(db *sql.DB, tableName string) func() *Query {
	return func() *Query {
		return &Query{
			db:    db,
			table: tableName,
		}
	}
}

func (q *Query) Insert(in interface{}) (int64, error) {
	var keys, values []string
	v := reflect.ValueOf(in)
	// 剥离指针
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		keys, values = sKV(v)
	case reflect.Map:
		keys, values = mKV(v)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++{
			//Kind是切片时，可以用Index()方法遍历
			sv := v.Index(i)
			for sv.Kind() == reflect.Ptr || sv.Kind() == reflect.Interface {
				sv = sv.Elem()
			}
			//切片元素不是struct或者指针，报错
			if sv.Kind() != reflect.Struct {
				return 0, errors.New("method Insert error: in slice is not structs")
			}
			//keys只保存一次就行，因为后面的都一样了
			if len(keys) == 0 {
				keys, values = sKV(sv)
				continue
			}
			_, val := sKV(sv)
			values = append(values, val...)
		}
	default:
		return 0, errors.New("method Insert error: type error")
		}

	}
	kl := len(keys)
	vl := len(values)
	if kl == 0 || vl == 0 {
		return 0, errors.New("method Insert error: no data")
	}
	var insertValue string
	//插入多条记录时需要用","拼接一下values
	if kl < vl {
		var tmpValues []string
		for kl <= vl {
			if kl%(len(keys)) == 0 {
				tmpValues = append(tmpValues, fmt.Sprintf("(%s)", strings.Join(values[kl-len(keys):kl], ",")))
			}
			kl++
		}
		insertValue = strings.Join(tmpValues, ",")
	} else {
		insertValue = fmt.Sprintf("(%s)", strings.Join(values, ","))
	}
	query := fmt.Sprintf(`insert into %s (%s) values %s`, q.table, strings.Join(keys, ","), insertValue)
	log.Printf("insert sql: %s", query)
	st, err := q.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	result, err := st.Exec()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()

}

func sKV(v reflect.Value) ([]string, []string) {
	var keys, values []string
	t := v.Type()
	for n := 0; n < t.NumField(); n++{
		tf := t.Field(n)
		vf := v.Field(n)

		if tf.Anonymous {
			continue
		}
		if !vf.IsValid() || reflect.DeepEqual(vf.Interface(), reflect.Zero(vf.Type()).Interface()){
			continue
		}
		for vf.Type().Kind() == reflect.Ptr {
			vf = vf.Elem()
		}

	//有时候根据需求会组合struct，这里处理下，支持获取嵌套的struct tag和value
		//如果字段值是time类型之外的struct，递归获取keys和values
		if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
			cKeys, cValues := sKV(vf)
			keys = append(keys, cKeys...)
			values = append(values, cValues...)
			continue
		}
		//根据字段的json tag获取key，忽略无tag字段
		key := strings.Split(tf.Tag.Get("json"), ",")[0]
		if key == "" {
			continue
		}
		value := format(vf)
		if value != "" {
			keys = append(keys, key)
			values = append(values, value)
		}
	}
	return keys, values

}

func format(v reflect.Value) string {
	//断言出time类型直接转unix时间戳
	if t, ok := v.Interface().(time.Time); ok {
		return fmt.Sprintf("FROM_UNIXTIME(%d)", t.Unix())
	}
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf(`'%s'`, v.Interface())
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return fmt.Sprintf(`%d`, v.Interface())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`%f`, v.Interface())
	//如果是切片类型，遍历元素，递归格式化成"(, , , )"形式
	case reflect.Slice:
		var values []string
		for i := 0; i < v.Len(); i++ {
			values = append(values, format(v.Index(i)))
		}
		return fmt.Sprintf(`(%s)`, strings.Join(values, ","))
	//接口类型剥一层递归
	case reflect.Interface:
		return format(v.Elem())
	}
	return ""
}

func mKV(v reflect.Value) ([]string, []string) {
	var keys, values []string
	//获取map的key组成的切片
	mapKeys := v.MapKeys()
	for _, key := range mapKeys {
		value := format(v.MapIndex(key))
		if value != "" {
			values = append(values, value)
			keys = append(keys, key.Interface().(string))
		}
	}
	return keys, values
}

// user1 := &User{
// 	Age:       30,
// 	FirstName: "Tom",
// 	LastName:  "Cat",
// }
// user2 := User{
// 	Age:       30,
// 	FirstName: "Tom",
// 	LastName:  "Curise",
// }
// user3 := User{
// 	Age:       30,
// 	FirstName: "Tom",
// 	LastName:  "Hanks",
// }
// user4 := map[string]interface{}{
// 	"age":        30,
// 	"first_name": "Tom",
// 	"last_name":  "Zzy",
// }
// users().Insert([]interface{}{user1, user2})
// users().Insert(user3)
// users().Insert(user4)

//Query will build a sql
type Query struct {
	db     *sql.DB
	table  string
	wheres []string
	only   []string
	limit  string
	offset string
	order  string
	errs   []string
}

//Where args can be string, User, *User, map[string]interface{}
func (q *Query) Where(wheres ...interface{}) *Query {
	for _, w := range wheres {
		str, err := where(true, w)
		q.wheres = append(q.wheres, str)
		if err != nil {
			//因为需要达到链式调用的效果，所以把错误都搜集起来，最后再处理
			q.errs = append(q.errs, err.Error())
		}
	}
	return q
}

func where(eq bool, w interface{}) (string, error) {
	var keys, values []string
	v := reflect.ValueOf(w)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return w.(string), nil
	case reflect.Struct:
		keys, values = sKV(v)
	case reflect.Map:
		keys, values = mKV(v)
	default:
		return "", errors.New("method Where error: type error")
	}
	if len(keys) != len(values) {
		return "", errors.New("method Where error: len(keys) not equal len(values))")
	}
	var wheres []string
	//之前的format()函数里，已经将切片类型值处理成"( , , ,)“形式
	for idx, key := range keys {
		if eq {
			if strings.HasPrefix(values[idx], "(") && strings.HasSuffix(values[idx], ")") {
				wheres = append(wheres, fmt.Sprintf("%s in %s", key, values[idx]))
				continue
			}
			wheres = append(wheres, fmt.Sprintf("%s = %s", key, values[idx]))
			continue
		}
		if strings.HasPrefix(values[idx], "(") && strings.HasSuffix(values[idx], ")") {
			wheres = append(wheres, fmt.Sprintf("%s not in %s", key, values[idx]))
			continue
		}
		wheres = append(wheres, fmt.Sprintf("%s != %s", key, values[idx]))
	}
	return strings.Join(wheres, " and "), nil
}

func (q *Query) Limit(limit uint) *Query {
	q.limit = fmt.Sprintf("limit %d", limit)
	return q
}

//Offset .
func (q *Query) Offset(offset uint) *Query {
	q.offset = fmt.Sprintf("offset %d", offset)
	return q
}

//Order .
func (q *Query) Order(ord string) *Query {
	q.order = fmt.Sprintf("order by %s", ord)
	return q
}

//Only 指定需要查询的字段
func (q *Query) Only(columns ...string) *Query {
	q.only = append(q.only, columns...)
	return q
}

func (q *Query) toSQL() string {
	var where string
	if len(q.wheres) > 0 {
		where = fmt.Sprintf(`where %s`, strings.Join(q.wheres, " and "))
	}
	sqlStr := fmt.Sprintf(`select %s from %s %s %s %s %s`, strings.Join(q.only, ","), q.table, where, q.order, q.limit, q.offset)
	log.Printf("select sql: %s", sqlStr)
	return sqlStr
}
//Select dest must be a ptr, e.g. *user, *[]user, *[]*user, *map, *[]map, *int, *[]int
func (q *Query) Select(dest interface{}) error {
	if len(q.errs) != 0 {
		return errors.New(strings.Join(q.errs, "
"))
	}
	t := reflect.TypeOf(dest)
	v := reflect.ValueOf(dest)
	typeErr := errors.New("method Select error: type error")
	if t.Kind() != reflect.Ptr {
		return typeErr
	}
	//如果是用 var userPtr *User 方式声明的变量，则不可取址
	if !v.Elem().CanAddr() {
		return typeErr
	}
	t = t.Elem()
	v = v.Elem()
	//如果only此时仍然为空，说明Only()方法未被调用，我们从struct上取tag填充
	if len(q.only) == 0 {
		switch t.Kind() {
		case reflect.Struct:
			if t.Name() != "Time" {
				q.only = sK(v)
			}
		case reflect.Slice:
			//获取切片的基本类型给一个局部变量
			t := t.Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() == reflect.Struct {
				if t.Name() != "Time" {
					q.only = sK(reflect.Zero(t))
				}
			}
		}
	}
	if len(q.only) == 0 {
		return errors.New("method Select error: type error, no columns to select")
	}
	if t.Kind() != reflect.Slice {
		q.limit = "limit 1"
	}

	rows, err := q.DB.Query(q.toSQL())
	if err != nil {
		return err
	}
	switch t.Kind() {
	case reflect.Slice:
		dt := t.Elem()
		for dt.Kind() == reflect.Ptr {
			dt = dt.Elem()
		}
		sl := reflect.MakeSlice(t, 0, 0)
		for rows.Next() {
			var destination reflect.Value
			if dt.Kind() == reflect.Map {
				destination, err = q.setMap(rows, dt)
			} else {
				destination, err = q.setElem(rows, dt)
			}
			if err != nil {
				return err
			}
			//区分切片元素是否指针
			switch t.Elem().Kind() {
			case reflect.Ptr, reflect.Map:
				sl = reflect.Append(sl, destination)
			default:
				sl = reflect.Append(sl, destination.Elem())
			}
		}
		v.Set(sl)
		return nil
	case reflect.Map:
		for rows.Next() {
			m, err := q.setMap(rows, t)
			if err != nil {
				return err
			}
			v.Set(m)
		}
		return nil
	default:
		for rows.Next() {
			destination, err := q.setElem(rows, t)
			if err != nil {
				return err
			}
			v.Set(destination.Elem())
		}
	}
	return nil

}

func sK(v reflect.Value) []string {
	var keys []string
	t := v.Type()
	for n := 0; n < t.NumField(); n++ {
		tf := t.Field(n)
		vf := v.Field(n)
		//忽略非导出字段
		if tf.Anonymous {
			continue
		}
		for vf.Type().Kind() == reflect.Ptr {
			vf = vf.Elem()
		}
		//如果字段值是time类型之外的struct，递归获取keys
		if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
			keys = append(keys, sK(vf)...)
			continue
		}
		//根据字段的json tag获取key，忽略无tag字段
		key := strings.Split(tf.Tag.Get("json"), ",")[0]
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}

func address(dest reflect.Value, columns []string) []interface{} {
	dest = dest.Elem()
	t := dest.Type()
	addrs := make([]interface{}, 0)
	switch t.Kind() {
	case reflect.Struct:
		for n := 0; n < t.NumField(); n++ {
			tf := t.Field(n)
			vf := dest.Field(n)
			if tf.Anonymous {
				continue
			}
			for vf.Type().Kind() == reflect.Ptr {
				vf = vf.Elem()
			}
			//如果字段值是time类型之外的struct，递归取址
			if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
				nVf := reflect.New(vf.Type())
				vf.Set(nVf.Elem())
				addrs = append(addrs, address(nVf, columns)...)
				continue
			}
			column := strings.Split(tf.Tag.Get("json"), ",")[0]
			if column == "" {
				continue
			}
			//只取选定的字段的地址
			for _, col := range columns {
				if col == column {
					addrs = append(addrs, vf.Addr().Interface())
					break
				}
			}
		}
	default:
		addrs = append(addrs, dest.Addr().Interface())
	}
	return addrs
}

//map的value类型必须是interface{}，因为无类型信息，所以mysql驱动会返回一个字节切片，需要自行用[]byte断言
func (q *Query) setMap(rows *sql.Rows, t reflect.Type) (reflect.Value, error) {
	if t.Elem().Kind() != reflect.Interface {
		return reflect.ValueOf(nil), errors.New("method setMap error: type error, must be map[string]interface{}")
	}
	m := reflect.MakeMap(t)
	addrs := make([]interface{}, len(q.only))
	for idx := range q.only {
		addrs[idx] = new(interface{})
	}
	if err := rows.Scan(addrs...); err != nil {
		return reflect.ValueOf(nil), err
	}
	for idx, column := range q.only {
		//从指针剥出interface{}，再剥出实际值
		m.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(addrs[idx]).Elem().Elem())
	}
	return m, nil
}
 //适用于基类型和struct
func (q *Query) setElem(rows *sql.Rows, t reflect.Type) (reflect.Value, error) {
	addrsErr := errors.New("method setElem error: columns not match addresses")
	dest := reflect.New(t)
	addrs := address(dest, q.only)
	if len(q.only) != len(addrs) {
		return reflect.ValueOf(nil), addrsErr
	}
	if err := rows.Scan(addrs...); err != nil {
		return reflect.ValueOf(nil), err
	}
	return dest, nil
}


//Update src can be *user, user, map[string]interface{}, string
func (q *Query) Update(src interface{}) (int64, error) {
	if len(q.errs) != 0 {
		return 0, errors.New(strings.Join(q.errs, "
"))
	}
	v := reflect.ValueOf(src)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var toBeUpdated, where string
	var keys, values []string
	switch v.Kind() {
	case reflect.String:
		toBeUpdated = src.(string)
	case reflect.Struct:
		keys, values = sKV(v)
	case reflect.Map:
		keys, values = mKV(v)
	default:
		return 0, errors.New("method Update error: type error")
	}
	if toBeUpdated == "" {
		if len(keys) != len(values) {
			return 0, errors.New("method Update error: keys not match values")
		}
		var kvs []string
		for idx, key := range keys {
			kvs = append(kvs, fmt.Sprintf("%s = %s", key, values[idx]))
		}
		toBeUpdated = strings.Join(kvs, ",")
	}
	if len(q.wheres) > 0 {
		where = fmt.Sprintf(`where %s`, strings.Join(q.wheres, " and "))
	}
	query := fmt.Sprintf("update %s set %s %s", q.table, toBeUpdated, where)
	st, err := q.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	result, err := st.Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//Delete no args
func (q *Query) Delete() (int64, error) {
	if len(q.errs) != 0 {
		return 0, errors.New(strings.Join(q.errs, "
"))
	}
	var where string
	if len(q.wheres) > 0 {
		where = fmt.Sprintf(`where %s`, strings.Join(q.wheres, " and "))
	}
	st, err := q.DB.Prepare(fmt.Sprintf(`delete from %s %s`, q.table, where))
	if err != nil {
		return 0, err
	}
	result, err := st.Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// http://www.imooc.com/article/details/id/76498