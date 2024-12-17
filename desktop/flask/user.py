from flask import Flask, request, render_template, flash, redirect
from forms import TodoForm
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from datetime import datetime
import time

app = Flask(__name__)
app.config['SECRET_KEY'] = '6e9cbf576a31bba80f0a34e35c2e678b1e2eba9885edaf99f3ce4aa2f5'
app.config["SQLALCHEMY_DATABASE_URI"] = "postgresql://mypostgres:dbelet@db:5432/mypostgres"
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
    todos = []
    for todo in TodoItem.query.order_by(TodoItem.data_completed.desc()).all():
        todo.id = str(todo.id)
        todos.append(todo)
    return render_template("index.html", title="Layout page", todos=todos)



@app.route("/add_todo", methods=['GET', 'POST'])
def add_todo():
    if request.method == 'POST':
        form = TodoForm() 
        todo_name = form.name.data
        todo_description = form.description.data
        completed = form.completed.data
        completed = True if completed in [True, 'True', 'true', 1] else False
        new_todo = TodoItem(
            name=todo_name,
            description=todo_description,
            completed=completed,
            data_completed=time.ctime()
        )
        db.session.add(new_todo)
        db.session.commit()
        flash("Todo successfully added!", "success")
        return redirect("/")
    else:
        form = TodoForm() 
    return render_template("add_todo.html", form=form)



@app.route('/update_todo/<int:id>', methods=['POST', 'GET'])
def update_todo(id):
    if request.method == 'POST':
        form = TodoForm(request.form)
        todo_name = form.name.data
        todo_description = form.description.data
        completed = form.completed.data
        completed = True if completed in [True, 'True', 'true', 1] else False
        db.session.query(TodoItem).filter(TodoItem.id == id).update({
            "name": todo_name,
            "description": todo_description,
            "completed": completed
        })
        db.session.commit()
        flash("Todo successfully updated!", "success")
        return redirect("/")

    else:
        form = TodoForm()
        todo = db.session.query(TodoItem).get(id)
        if todo:
            form.name.data = todo.name
            form.description.data = todo.description
            form.completed.data = todo.completed
    return render_template("add_todo.html", form=form)


@app.route('/update_all_todos/<int:status>', methods=['GET'])
def update_all_todos(status):
    if status == 1: 
        todos = db.session.query(TodoItem).filter_by(completed=True).all()
    elif status == 0: 
        todos = db.session.query(TodoItem).filter_by(completed=False).all()
    else:  
        flash("Geçersiz durum!", "danger")
        return redirect("/")
    return render_template("index.html", todos=todos)



@app.route('/delete_todo/<int:id>', methods=['POST', 'GET'])
def delete_todo(id):
    todo = db.session.query(TodoItem).get(id)
    db.session.delete(todo)
    db.session.commit()
    flash("Todo deleted", "success")
    return redirect("/")

  
if __name__== '__main__':
    app.run(debug=True)

