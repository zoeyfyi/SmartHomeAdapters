package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.disposables.Disposable
import io.reactivex.schedulers.Schedulers
import kotlinx.android.synthetic.main.activity_authentication.*
import org.jetbrains.anko.*
import org.jetbrains.anko.design.snackbar
import org.json.JSONException
import org.json.JSONObject
import retrofit2.HttpException

class AuthenticationActivity : AppCompatActivity() {

    private val tag = "AuthenticationActivity"

    private val restApiService by lazy {
        RestApiService.create()
    }

    private var disposable: Disposable? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        /*
            Set up appropriate listeners for button click events
         */
        sign_in_button.setOnClickListener { _ ->

            sign_in_button.isEnabled = false
            login_progress_bar.visibility = View.VISIBLE

            val email = email_input.text.toString()
            val password = password_input.text.toString()
            val inputsOK = checkEmailInput(email) && checkPasswordInput(password)

            if (inputsOK) {
                signInUser(User(email, password))
            }
        }

        register_button.setOnClickListener { _ ->
            Log.d(tag, "Starting RegisterUserActivity")
            startActivity<RegisterUserActivity>()
        }

        // Validate the inputs whenever they lose focus
        email_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                // The email input has dropped focus; check its validity.
                checkEmailInput(email_input.text.toString())
            }
        }
        password_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                // The password input has dropped focus; check its validity.
                checkPasswordInput(password_input.text.toString())
            }
        }
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
        login_progress_bar.visibility = View.GONE
        sign_in_button.isEnabled = true
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

    override fun onPause() {
        super.onPause()
        disposable?.dispose()
    }
}
