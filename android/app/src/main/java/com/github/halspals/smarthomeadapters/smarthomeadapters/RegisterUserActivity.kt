package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.disposables.Disposable
import io.reactivex.schedulers.Schedulers
import kotlinx.android.synthetic.main.activity_register_user.*
import org.jetbrains.anko.*
import org.jetbrains.anko.design.snackbar
import org.json.JSONException
import org.json.JSONObject
import retrofit2.HttpException
import java.net.HttpURLConnection

class RegisterUserActivity : AppCompatActivity() {

    private val tag = "RegisterUserActivity"

    private lateinit var user: User

    private val restApiService by lazy {
        RestApiService.create()
    }

    private var disposable: Disposable? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_user)

        /*
            Set up the click events for the buttons.
         */
        register_button.setOnClickListener { _ ->
            register_button.isEnabled = false
            progressBar.visibility = View.VISIBLE

            val email = email_input.text.toString()
            val password = password_input.text.toString()
            val confirmedPassword = confirm_password_input.text.toString()

            val inputsOK = checkEmailInput(email)
                    && checkPasswordInput(password)
                    && checkConfirmPasswordInput(password, confirmedPassword)

            if (inputsOK) {
                registerNewUser(User(email, password))
            } else {
                register_button.isEnabled = true
                progressBar.visibility = View.GONE
            }
        }


        login_button.setOnClickListener { _ ->
            // Return to the AuthenticationActivity which launched this RegisterUserActivity
            finish()
        }

        /*
            Set onFocusChangeListeners for the input fields to check if their input thus far
            is valid.
         */
        email_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                checkEmailInput(email_input.text.toString())
            }
        }

        password_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                checkPasswordInput(password_input.text.toString())
            }
        }

        confirm_password_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                checkConfirmPasswordInput(
                        password_input.text.toString(), confirm_password_input.text.toString())
            }
        }


    }

    /**
     * Gets the input email and password and registers a new user with the web service.
     * First makes sure that the inputs are valid.
     *
     */
    private fun registerNewUser(user: User) {
        disposable = restApiService.registerUser(user)
                .subscribeOn(Schedulers.io())
                .observeOn(AndroidSchedulers.mainThread())
                .subscribe(
                        { user_confirmation ->
                            assert(user == user_confirmation) { "Returned user doesn't equal expected user. " +
                                "Expected {$user}, received {$user_confirmation}" }
                            toast("Registered; now signing you in...")
                            signInUser(user) },
                        { error -> handleLoginError(error) }
                )
    }

    /**
     * WIP: Get the email and password given, validate them, and then
     * try to authenticate the user.
     *
     * @return whether the inputs are valid and the authentication was successful
     */
    private fun signInUser(user: User) {
        disposable = restApiService.loginUser(user)
                .subscribeOn(Schedulers.io())
                .observeOn(AndroidSchedulers.mainThread())
                .subscribe(
                        { token -> saveTokenAndMoveToMain(token.token) },
                        { error -> handleLoginError(error) }
                )
    }

    /**
     * WIP: Saves the token and starts a [MainActivity].
     * Currently only keeps token in memory, passing it as an extra to the activity;
     * when we have time this should be changed to saving the token in Account Manager.
     *
     * @param token the authorization token received from the server
     */
    private fun saveTokenAndMoveToMain(token: String) {
        Log.d(tag, "Succeeded in receiving token, starting MainActivity")
        // TODO token should be saved in Account Manager for the user
        toast("Signed in")
        startActivity(intentFor<MainActivity>("token" to token).clearTask().newTask())
    }

    /**
     * Handles an error received by [signInUser], displaying the message to the user.
     *
     * @param error the error received from the api call
     */
    private fun handleLoginError(error: Throwable) {
        // There was an error; if the server gave an error message in JSON format,
        // try to extract it
        val errorString: String? = if (error is HttpException) {
            try {
                JSONObject(error.response().errorBody()?.string()).getString("error")
            } catch (e: JSONException) {
                error.message.toString()
            }
        } else {
            // If the error is not an HttpException we will have to make to
            // with its error message
            error.message.toString()
        }

        // Display the error to the user
        Log.d(tag, "Login failed: $errorString")
        if (errorString != null) {
            snackbar_layout.snackbar(errorString)
        }

        // Allow the user to try again
        register_button.isEnabled = true
        progressBar.visibility = View.GONE
    }
    /**
     * Makes sure the given email is valid, setting an error to [email_input] if not.
     *
     * @param email the email to check for validity
     * @return whether the given email is a valid one
     */
    private fun checkEmailInput(email: String): Boolean {

        val emailIsValid = android.util.Patterns.EMAIL_ADDRESS.matcher(email).matches()

        if (!emailIsValid) {
            email_input.error = "Invalid email address"
        } else {
            email_input.error = null
        }

        return emailIsValid
    }


    /**
     * Checks the validity of the password (>= 8 chars).
     * Sets an error to [password_input] if it is not valid.
     *
     * @param password the password to check
     * @return whether the password is a valid one
     */
    private fun checkPasswordInput(password: String): Boolean {
        // TODO we can add more checks here like password length etc
        val pwIsValid = password.length >= 8

        if (!pwIsValid) {
            password_input.error = "Must be at least 8 characters long"
        } else {
            password_input.error = null
        }

        return pwIsValid
    }

    /**
     * Checks that the two passwords match, adding an error to [confirm_password_input] if not.
     *
     * @return whether the two passwords match
     */
    private fun checkConfirmPasswordInput(password: String, confirmed_password: String): Boolean {
        // Simply make sure that the two passwords given match.

        val confirmedPwIsValid = password == confirmed_password

        if (!confirmedPwIsValid) {
            confirm_password_input.error = "Passwords do not match"
        } else {
            confirm_password_input.error = null
        }

        return confirmedPwIsValid
    }

    override fun onPause() {
        super.onPause()
        disposable?.dispose()
    }
}
