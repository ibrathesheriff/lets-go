package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.ibrathesheriff.com/internal/models"
	"snippetbox.ibrathesheriff.com/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator  `form:"-"`
}

// Create  a new userSignupForm struct.
type userSignupForm struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`

}

// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the
	// manual check of r.URL.Path != "/" from this handler.

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

// Add a new snippetCreate handler, which for now returns a placeholder
// response. We'll update this shortly to show a HTML form.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Declare a new empty instance of the snippetCreateForm struct.
	var form snippetCreateForm

	// Call the Decode() method of the form decoder, passing in the current
	// request and *a pointer* to our snippetCreateForm struct. This will
	// essentially fill our struct with the relevant values from the HTML form.
	// If there is a problem, we return a 400 Bad Request response to the client.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Then validate and use the data as normal...
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the Valid() method to see if any of the checks failed. If they did,
	// then re-render the template passing in the form in the same way as
	// before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Snippet successfully
	// created!") and the corresponding key ("flash") to the session data.
	// - The first parameter that we pass to app.sessionManager.Put() is the
	//   current request context
	// - The second parameter - "flash" -  is the key for the specific message
	//   that we are adding to the session data
	// - If there's no existing session for the current user (or their session
	//   has expired) then a new, empty, session for them will automatically be
	//   created by the session middleware.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Update the handler so it displays the signup page.

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm

	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using our helper functions
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

