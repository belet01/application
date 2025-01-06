from flask import Flask, request, render_template, flash, redirect, session, url_for, jsonify
from forms import TodoForm, RegisterForm, LoginForm
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
import time

app = Flask(__name__)
app.config['SECRET_KEY'] = '6e9cbf576a31bba80f0a34e35c2e678b1e2eba9885edaf99f3ce4aa2f5'
app.config["SQLALCHEMY_DATABASE_URI"] = "postgresql://mypostgres:dbelet@172.18.0.2/dreamer"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
db = SQLAlchemy(app)
migrate = Migrate(app, db)


class Users(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), nullable=False, unique=True)
    email = db.Column(db.String(120), nullable=False, unique=True)
    password_hash = db.Column(db.String(900), nullable=False)
    todos = db.relationship('TodoItem', backref='user', lazy=True)
    bio = db.Column(db.String(300), nullable=True)
    name = db.Column(db.String(80), nullable=True)
    followers = db.Table('followers',
        db.Column('follower_id', db.Integer, db.ForeignKey('users.id'), primary_key=True),
        db.Column('followed_id', db.Integer, db.ForeignKey('users.id'), primary_key=True)
    )
    followed = db.relationship('Users', secondary=followers,
                               primaryjoin=(id == followers.c.follower_id),
                               secondaryjoin=(id == followers.c.followed_id),
                               backref=db.backref('followers', lazy='dynamic'),
                               lazy='dynamic')


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
    theme = session.get('theme', 'creamy')  # Kullanıcının teması veya varsayılan 'creamy'
    followed_users = []  # Takip edilen kullanıcılar

    if 'username' in session:
        username = session['username']
        me = Users.query.filter_by(username=username).first()
        followed_users = me.followed.all()  # Kullanıcının takip ettiği kişiler
        todos = TodoItem.query.order_by(TodoItem.data_completed.desc()).all()
        return render_template("index.html", title="Layout page", todos=todos, me=me, theme=theme, followed_users=followed_users)

    todos = TodoItem.query.order_by(TodoItem.data_completed.desc()).all()
    return render_template("index.html", title="Layout page", todos=todos, theme=theme)




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
    if request.method== 'POST':
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
    user = Users.query.filter_by(username=username).first_or_404()
    todos = TodoItem.query.filter_by(user_id=user.id).all()
    todo_count = len(todos)
    follower_count = user.followers.count()  
    followee_count = user.followed.count()
    current_user = session.get('username')
    followed_users = []
    if current_user:
        current_user_obj = Users.query.filter_by(username=current_user).first()
        followed_users = current_user_obj.followed.all()
    if current_user == username:
        return render_template('myprofil.html', user=user, followers_count=follower_count, 
                               followees_count=followee_count, todo_count=todo_count, todos=todos)
    return render_template('user_profile.html', user=user, followers_count=follower_count, 
                           followees_count=followee_count, todo_count=todo_count, todos=todos, 
                           followed_users=followed_users)



@app.route('/my_profile/<username>')
def my_profile(username):
    user = Users.query.filter_by(username=username).first()
    todos = TodoItem.query.filter_by(user_id=user.id).all()
    todo_count = len(todos)
    follower_count = user.followers.count()  
    followee_count = user.followed.count()
    return render_template('myprofil.html', user=user, followers_count=follower_count, followees_count=followee_count, todo_count= todo_count, todos= todos)



@app.route("/profil/update/<int:id>", methods=['GET', 'POST'])
def profile_settings(id):
    users = db.session.query(Users).get(id)
    if request.method == 'POST':
        form = LoginForm(request.form)
        username = form.username.data
        password = form.password.data
        name = form.name.data  
        biografi = form.biografi.data
        if Users.query.filter(Users.username == username, Users.id != id).first():
            flash("Username already exists!", "danger")
            return redirect(url_for('profile_settings', id=users.id))
        else: 
            db.session.query(Users).filter(Users.id == id).update({
                Users.username: username,
                Users.password_hash: password,
                Users.name: name,  # Name alanını güncelliyoruz
                Users.bio: biografi  # Biografi alanını güncelliyoruz
            })
            db.session.commit()

        # Session bilgilerini güncelleniyor
        session['username'] = username
        session['password_hash'] = password
        session['name'] = name  # Name bilgisi session'a ekleniyor
        session['bio'] = biografi  # Biyografi bilgisi session'a ekleniyor
        flash("Profil başarıyla güncellendi!", "success")
        return redirect(url_for('home', username=session['username']))
    else:
        form = LoginForm()
        if users:
            form.username.data = users.username
            form.password.data = users.password_hash
            form.name.data = users.name if users.name else ""  # Name form'a ekleniyor
            form.biografi.data = users.bio if users.bio else ""  # Biografi form'a ekleniyor

    return render_template("profil_settings.html", form=form, users=users)




@app.route('/change_theme', methods=['POST'])
def change_theme():
    data = request.get_json()
    theme = data.get('theme')
    if theme in ['light', 'dark', 'creamy']:
        session['theme'] = theme
        return jsonify({'success': True})
    return jsonify({'success': False})



@app.route('/change_theme_page')
def change_theme_page():
    return render_template('theme_settings.html')

@app.route("/all/users")
def all_users():
    usernames = db.session.query(Users.username).all()
    title = "All users:"
    return render_template("users.html", usernames= usernames, title= title)


@app.route('/follow/<int:user_id>', methods=['POST'])
def follow(user_id):
    user_to_follow = Users.query.get_or_404(user_id)
    username = session.get('username')
    current_user = Users.query.filter_by(username=username).first()
    if current_user.id == user_to_follow.id:
        flash("You cannot follow yourself.", "danger")
        return redirect(url_for('profile', username=current_user.username))
        
    if user_to_follow not in current_user.followed:
        current_user.followed.append(user_to_follow)
        db.session.commit()
        flash("You have started following this user.", "success")
    else:
        current_user.followed.remove(user_to_follow)
        db.session.commit()
        flash("You have unfollowed this user.", "success")
    return redirect(url_for('home', username=user_to_follow.username ))




@app.route('/<username>/connections/<string:connection_type>')
def user_connections(username, connection_type):
    user = Users.query.filter_by(username=username).first_or_404()
    if connection_type == "followed":  # Takip ettikleri
        connections = user.followed.all()
        title = f"{user.username}'s follow"
    elif connection_type == "followers":  # Takip edenler
        connections = user.followers.all()
        title = f"{user.username}'s follower"
    else:
        title = "Hata"
    return render_template("connections.html", user=user, connections=connections, connection_type=connection_type, title=title)


@app.route('/username/delete/<id>',  methods=['POST', 'GET'])
def delete_account(id):
    user= db.session.query(Users).get(id)
    todos = TodoItem.query.filter_by(user_id=id).all()
    for todo in todos:
        db.session.delete(todo)
    db.session.delete(user)
    db.session.commit()
    session.pop('username', None)
    return redirect("/register")

if __name__== '__main__':
    app.run(debug=True)
