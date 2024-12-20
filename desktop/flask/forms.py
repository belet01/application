from flask_wtf import FlaskForm
from wtforms import StringField, TextAreaField, SelectField, SubmitField, PasswordField
from wtforms.validators import DataRequired, EqualTo, Length


class TodoForm(FlaskForm):
    name = StringField('Name', validators=[DataRequired()])
    description = TextAreaField('Description', validators=[DataRequired()])
    completed = SelectField('Completed',choices=[(False, 'False'),(True, 'True') ],validators=[DataRequired()])
    submit = SubmitField('Add Todo')

class RegisterForm(FlaskForm):
    username = StringField('Username', validators=[DataRequired(), Length(min=3, max=80)])
    email = StringField('Email', validators=[DataRequired()])
    password = PasswordField('Password', validators=[DataRequired(), Length(min=6)])
    confirm_password = PasswordField('Confirm Password', validators=[DataRequired(), EqualTo('password')])
    submit = SubmitField('Register')

class LoginForm(FlaskForm):
    username = StringField('Username', validators=[DataRequired()])
    password = PasswordField('Password', validators=[DataRequired(), Length(min=6)])