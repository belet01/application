from flask import Flask, request
import jwt
import psycopg2
from dream import get_token, get_user_id_from_token
app = Flask(__name__)

@app.route("/students")
def home():
    conn= psycopg2.connect(host= '172.18.0.2', database = 'postgres', user = 'mypostgres', password = 'dbelet')
    cur = conn.cursor()
    cur.execute("SELECT * FROM student")
    data = cur.fetchall()
    for do in data:
        return {
        'students': [
            {
            "id": do[0], 
            "name": do[1],
            "username": do[2],
            "surname": do[3],
            "age": do[4],
            "password": do[5] 
            }
        ]
    }


@app.route("/register", methods = ['POST'])
def register():
    conn= psycopg2.connect(host= '172.18.0.2', database = 'mypostgres', user = 'mypostgres', password = 'dbelet')
    cur = conn.cursor()
    cur.execute('SELECT username FROM student')
    data_base = cur.fetchall()
    for username in data_base:
        if username[0] == request.json['username']:
                return {'message' : 'Username exits!'}
    else:              
        cur.execute('INSERT INTO student(name, surname, age, username, password) VALUES( %s, %s, %s, %s, %s)',
                (request.json['name'], request.json['surname'], request.json['age'], request.json['username'], request.json['password']))
        conn.commit()
        cur.close()
        conn.close()            
        return {'status' : 'successful'}
    

@app.route("/login", methods = ['POST'])
def login():
    user_name = request.json['username']
    password = request.json['password']
    conn= psycopg2.connect(host= '172.18.0.2', database = 'mypostgres', user = 'mypostgres', password = 'dbelet')
    cur = conn.cursor()
    cur.execute("SELECT id, name, surname, username, password, age FROM student WHERE password=%s and username =%s", (password, user_name))
    data = cur.fetchall()
    for do in data:
        if do[3] == user_name and do[4] == password:
             return  {
                    "token": get_token(int(do[0])),
                    "name": do[1],
                    "surname" : do[2],
                    "age" : do[5]
                }
    else:
        return {'message' : 'Your username or password is not found... Try again!'}
   

@app.route("/put", methods = ['PUT'])
def put():
    try: 
        token = request.headers['Authorization']
        id = get_user_id_from_token(token)
    except jwt.ExpiredSignatureError as err:
        return {'messege': f"{err}"}
    except jwt.InvalidTokenError as err:
        return {'messege': f"{err}"}
    else:
        conn= psycopg2.connectconn= psycopg2.connect(host= '172.18.0.2', database = 'mypostgres', user = 'mypostgres', password = 'dbelet')
        cur = conn.cursor()
        cur.execute('SELECT id FROM student')
        data = cur.fetchall()
        for user_id in data:
            if user_id[0] == id:
                cur.execute("SELECT id, username FROM student")
                data1 = cur.fetchall()
                for do in data1:
                    if do[1] == request.json['username'] and do[0] != id: 
                        return {"message" : "Username exits"}
                else:
                    cur.execute("UPDATE student SET name=%s, username=%s, surname=%s, age=%s, password=%s WHERE id =%s", 
                        ( request.json['name'], request.json['username'], request.json['surname'], request.json['age'], request.json['password'], id) )
                    conn.commit()
                    cur.close()
                    conn.close()
                    return {'message' : 'Successful'}


@app.route('/delete', methods = ['DELETE'])
def delete():
    try: 
        token = request.headers['Authorization']
        id = get_user_id_from_token(token)
    except jwt.ExpiredSignatureError as err:
        return {'messege': f"{err}"}
    except jwt.InvalidTokenError as err:
        return {'messege': f"{err}"}
    else:
        conn= psycopg2.connect(host= '172.18.0.2', database = 'mypostgres', user = 'mypostgres', password = 'dbelet')
        cur = conn.cursor()
        cur.execute('SELECT id FROM student')
        data = cur.fetchall()
        for user_id in data:
            if user_id[0] == id:
                cur.execute('DELETE FROM student WHERE id = %s', (user_id))
                cur.close()
                conn.commit()
                conn.close()
                return {"message": "deleted success!"}


@app.route("/patch", methods = ['PATCH'])
def patch():
    try: 
        token = request.headers['Authorization']
        id = get_user_id_from_token(token)
    except jwt.ExpiredSignatureError as err:
        return {'messege': f"{err}"}
    except jwt.InvalidTokenError as err:
        return {'messege': f"{err}"}
    else:
        ada = list(request.get_json().keys())
        if len(ada) == 0:
            return {"message":"you can't send it blank"}
        else:
            conn= psycopg2.connect(host= '172.18.0.2', database = 'mypostgres', user = 'mypostgres', password = 'dbelet')
            cur = conn.cursor()
            cur.execute('SELECT id FROM student')
            data = cur.fetchall()
            for user_id in data:
                if user_id[0] == id:
                    if ada[0] == 'username':
                        cur.execute("SELECT id, username FROM student")
                        database = cur.fetchall()
                        for do in database:
                            if  do[0] != id and do[1] == request.json['username']: 
                                return {"message" : "Username exits"}
                        else:
                            cur.execute("UPDATE student SET username=%s WHERE id =%s", (request.json['username'], id))
                            cur.close()
                            conn.commit()
                            cur.close()
                            return {'message' : 'Successful'}
                    else:
                        liste = ['name', 'surname', 'age', 'password']
                        if ada[0] not in liste:
                            return {'message': 'Wrong key!'}
                        else:
                            cur.execute("UPDATE student SET {}=%s WHERE id =%s".format(ada[0]),(request.json[ada[0]], id) )
                            cur.close()
                            conn.commit()
                            cur.close()
                            return {'message' : 'Successful'}

@app.route('/student-info')
def user_info():
    try:
        token = request.headers['Authorization']
        id = get_user_id_from_token(token)
    except jwt.ExpiredSignatureError as err:
        return {'messege': f"{err}"}
    except jwt.InvalidTokenError as err:
        return {'messege': f"{err}"}
    else:
        conn= conn= psycopg2.connect(host= '172.18.0.2', database = 'mypostgres', user = 'mypostgres', password = 'dbelet')
        cur = conn.cursor()
        cur.execute('SELECT id FROM student')
        data = cur.fetchall()
        for user_id in data:
            if user_id[0] == id:
                cur.execute("SELECT name, surname, username, age FROM student WHERE id = '{}' ".format(id))
                data = cur.fetchall()
                for do in data:
                    return {
                        "name": do[0],
                        "surname" : do[1],
                        "username" : do[2], 
                        "age": do[3]
                        }

    
if __name__== '__main__':
    app.run(debug=True)
