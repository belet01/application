from flask import Flask, request, render_template, flash, redirect, session, url_for
from forms import TodoForm, RegisterForm, LoginForm
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
import time


app = Flask(__name__)
app.config['SECRET_KEY'] = '6e9cbf576a31bba80f0a34e35c2e678b1e2eba9885edaf99f3ce4aa2f5'
app.config["SQLALCHEMY_DATABASE_URI"] = "postgresql://mypostgres:dbelet@172.18.0.2/dreams"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
db = SQLAlchemy(app)
migrate = Migrate(app, db)


class Users(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), nullable=False, unique=True)
    email = db.Column(db.String(120), nullable=False, unique=True)
    password_hash = db.Column(db.String(900), nullable=False)
    todos = db.relationship('TodoItem', backref='user', lazy=True)


class TodoItem(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    description = db.Column(db.String(200), nullable=True)
    completed = db.Column(db.Boolean, default=False)
    data_completed = db.Column(db.DateTime, nullable=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=True)


with app.app_context():
    db.create_all()

@app.route('/update_all_todos/<int:status>', methods=['GET'])
def update_all_todos(status):
    if status == 1: 
        todos = db.session.query(TodoItem).filter_by(completed=True).all()
    elif status == 0: 
        todos = db.session.query(TodoItem).filter_by(completed=False).all() 
    return render_template("index.html", todos=todos)


@app.route("/")
def home():
    if 'username' in session:
        username = session['username']
        me = Users.query.filter_by(username=username).first()
        todos = TodoItem.query.order_by(TodoItem.data_completed.desc()).all()
        return render_template("index.html", title="Layout page", todos=todos, me=me)
    todos = TodoItem.query.order_by(TodoItem.data_completed.desc()).all()
    return render_template("index.html", title="Layout page", todos=todos)




@app.route("/register", methods=["GET", "POST"])
def register():
    if 'username' in session:
        return redirect(url_for('home'))
    
    form = RegisterForm()
    if request.method == 'POST' and form.validate_on_submit():
        username = form.username.data
        email = form.email.data
        password = form.password.data

        if Users.query.filter((Users.username == username) | (Users.email == email)).first():
            flash("Username or email already exists!", "danger")
            return redirect("/register")
        else:
            new_user = Users(username=username, email=email, password_hash=password)
            db.session.add(new_user)
            db.session.commit()
            flash("Registration successful!", "success")
            return redirect("/login")
    return render_template("register.html", form=form)






@app.route("/login", methods=["GET", "POST"])
def login():
    if 'username' in session:
        return redirect(url_for('home'))
    
    form = LoginForm()
    if request.method == 'POST':
        username = form.username.data
        password = form.password.data
        search = Users.query.filter_by(username=username, password_hash=password).first()
        if search is None:
            flash("Kullanici adi ve ya sifre yalnis!", "danger")
            return render_template('login.html', form=form)
        
        session['username'] = username 
        return redirect(url_for('home'))
    return render_template("login.html", form=form)

@app.route("/add_todo", methods=['GET', 'POST'])
def add_todo():
    if 'username' not in session:
        flash("Todo ekleyemezsiniz!", "Hata")
        return redirect(url_for('login')) 
    form = TodoForm() 
    if request.method== 'POST' and form.validate_on_submit():
        todo_name = form.name.data
        todo_description = form.description.data
        completed = form.completed.data
        completed = True if completed in [True, 'True', 'true', 1] else False
        new_todo = TodoItem(
            name=todo_name,
            description=todo_description,
            completed=completed,
            data_completed=time.ctime(),
            user_id=Users.query.filter_by(username=session['username']).first().id
        )
        db.session.add(new_todo)
        db.session.commit()
        flash("Todo successfully added!", "success")
        return redirect(url_for('home'))
    
    return render_template("add_todo.html", form=form)


@app.route("/update_todo/<int:id>", methods=['GET', 'POST'])
def update_todo(id):
    if 'username' not in session:
        flash("You must be logged in to update todos.", "danger")
        return redirect(url_for('login'))
    todo = db.session.query(TodoItem).get(id)
    if todo.user_id != Users.query.filter_by(username=session['username']).first().id:
        flash("You cannot update someone else's todo.", "danger")
        return redirect(url_for('home'))
    if request.method == 'POST':
        form = TodoForm(request.form)
        todo_name = form.name.data
        todo_description = form.description.data
        completed = form.completed.data
        completed = True if completed in [True, 'True', 'true', 1] else False
        db.session.query(TodoItem).filter(TodoItem.id == id).update({
            "name": todo_name,
            "description": todo_description,
            "completed": completed,
        })
        db.session.commit()
        flash("Todo successfully updated!", "success")
        return redirect(url_for('my_profile', username=session['username']))

    else:
        form = TodoForm()
        if todo:
            form.name.data = todo.name
            form.description.data = todo.description
            form.completed.data = todo.completed
    return render_template("add_todo.html", form=form)


@app.route('/delete_todo/<int:id>', methods=['POST', 'GET'])
def delete_todo(id):
    if 'username' not in session:
        flash("You must be logged in to delete todos.", "danger")
        return redirect(url_for('login'))
    todo = db.session.query(TodoItem).get(id)
    if todo.user_id != Users.query.filter_by(username=session['username']).first().id:
        flash("You cannot delete someone else's todo.", "danger")
        return redirect(url_for('home'))
    db.session.delete(todo)
    db.session.commit()
    flash("Todo deleted", "success")
    return redirect("/")



@app.route('/logout', methods=['GET', 'POST'])
def logout():
    if request.method == 'POST':
        session.pop('username', None)
        flash("You have been logged out.", "success")
        return redirect(url_for('login'))
    return render_template('logout.html')


@app.route('/user_profile/<username>')
def profile(username):
    current_user = session.get('username')
    user = Users.query.filter_by(username=username).first()
    todos = TodoItem.query.filter_by(user_id=user.id).all()
    todo_count = len(todos)
    if current_user == username:
        return render_template('myprofil.html', user=user, todos=todos, todo_count=todo_count)
    return render_template('user_profile.html', user=user, todos=todos, todo_count=todo_count)


@app.route('/my_profile/<username>')
def my_profile(username):
    user = Users.query.filter_by(username=username).first()
    todos = TodoItem.query.filter_by(user_id=user.id).all()
    todo_count = len(todos)
    return render_template('myprofil.html', user=user, todos=todos, todo_count=todo_count)


@app.route("/profil/update/<int:id>", methods=['GET', 'POST'])
def profile_settings(id):
    users = db.session.query(Users).get(id)
    if request.method == 'POST':
        form = LoginForm(request.form)
        username = form.username.data
        password = form.password.data
        
        db.session.query(Users).filter(Users.id == id).update({
            "username": username,
            "password_hash": password  
        })
        db.session.commit()
        session['username'] = username
        session['password_hash'] = password 
        flash("Profil başarıyla güncellendi!", "success")
        return redirect(url_for('my_profile', username=session['username']))
    else:
        form = LoginForm()
        if users:
            form.username.data = users.username
            form.password.data = users.password_hash 

    return render_template("profil_settings.html", form=form, users=users)


@app.route('/theme/settings', methods=['GET', 'POST'])
def theme_settings():
    if request.method == 'POST':
        selected_theme = request.form.get('theme')  
        session['theme'] = selected_theme  
        return redirect(url_for('home'))  
    return render_template('theme_settings.html')



@app.route('/username/delete/<id>',  methods=['POST', 'GET'])
def delete_account(id):
    user= db.session.query(Users).get(id)
    db.session.delete(user)
    db.session.commit()
    session.pop('username', None)
    return redirect("/login")

if __name__== '__main__':
    app.run(debug=True)
