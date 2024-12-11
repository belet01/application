from flask import Flask, request, render_template, flash, redirect
from forms import TodoFrom
from flask_sqlalchemy import SQLAlchemy
import psycopg2
from flask_migrate import Migrate
from datetime import datetime

app = Flask(__name__)
app.config['SECRET_KEY'] = '6e9cbf576a31bba80f0a34e35c2e678b1e2eba9885edaf99f3ce4aa2f5'
app.config["SQLALCHEMY_DATABASE_URI"] = "postgresql://mypostgres:dbelet@172.18.0.2/mypostgres"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
db = SQLAlchemy(app)
migrate = Migrate(app, db)


class TodoItem(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    description = db.Column(db.String(200), nullable=True)
    completed = db.Column(db.Boolean, default=False)
    data_completed = db.Column(db.DateTime, nullable=True)


with app.app_context():
    db.create_all()


@app.route("/")
def home():
    return render_template("index.html")



@app.route("/add_todo", methods=['GET', 'POST'])
def add_todo():
    return render_template("add_todo.html")



@app.route('/update', methods = ['POST', 'GET'])
def delete():
    return " "

@app.route('/delete', methods = ['DELETE'])
def delete():
    return " "


    
if __name__== '__main__':
    app.run(debug=True)

