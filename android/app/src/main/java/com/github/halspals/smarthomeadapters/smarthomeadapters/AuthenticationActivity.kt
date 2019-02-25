package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.content.Intent
import android.net.Uri
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.view.inputmethod.InputMethodManager
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Token
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import kotlinx.android.synthetic.main.activity_authentication.*
import net.openid.appauth.*
import org.jetbrains.anko.*
import org.jetbrains.anko.design.snackbar
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class AuthenticationActivity : AppCompatActivity() {

    private val tag = "AuthenticationActivity"

    private val restApiService by lazy {
        RestApiService.new()
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        /*
            Set up appropriate listeners for button click events
         */
        sign_in_button.setOnClickListener { _ ->

            // Dismiss the keyboard
            val inputMethodManager = getSystemService(Context.INPUT_METHOD_SERVICE) as InputMethodManager
            val currentView = currentFocus ?: View(this)
            inputMethodManager.hideSoftInputFromWindow(currentView.windowToken, 0)

            // Signal to the user that we are waiting for an async call
            sign_in_button.isEnabled = false
            login_progress_bar.visibility = View.VISIBLE

            val email = email_input.text.toString()
            val password = password_input.text.toString()
            val inputsOK = checkEmailInput(email) && checkPasswordInput(password)

            if (inputsOK) {
                signInUser(User(email, password))
            } else {
                sign_in_button.isEnabled = true
                login_progress_bar.visibility = View.GONE
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

        // Set up Authorization configuration
        val authServiceConfig = AuthorizationServiceConfiguration(
            Uri.parse("https://oauth.halspals.co.uk/oauth2/auth"),
            Uri.parse("https://oauth.halspals.co.uk/oauth2/token")
        )

        // Get an authorization code
        val authRequest: AuthorizationRequest = AuthorizationRequest.Builder(
            authServiceConfig,
            "b43ce28c-f4e3-412b-8dc5-854a32a0c8db",
            ResponseTypeValues.CODE,
            Uri.parse("http://callback.halspals.co.uk")
        ).build()

        // Do the authorization
        Log.d(tag, "Starting Oauth2 call")
        val authService = AuthorizationService(this)
        val authIntent = authService.getAuthorizationRequestIntent(authRequest)
        startActivityForResult(authIntent, 42)
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        if (requestCode == 42) {
            Log.d(tag, "Received result from Oauth2 call")
            //val resp = AuthorizationResponse.fromIntent(data!!)
        } else {
            Log.w(tag, "Received UNEXPECTED activity result")
        }
    }

    /**
     * Makes a REST call to sign in the given user.
     *
     * @param user the user to sign in
     */
    private fun signInUser(user: User) {
        restApiService.loginUser(user).enqueue(object: Callback<Token> {
            override fun onResponse(call: Call<Token>, response: Response<Token>) {
                val token = response.body()
                if (response.isSuccessful && token != null) {
                    saveTokenAndMoveToMain(token.token)
                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    handleLoginError(error)
                }
            }

            override fun onFailure(call: Call<Token>, error: Throwable) {
                handleLoginError(error.message)
            }
        })
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
    private fun handleLoginError(error: String?) {

        // Display the error to the user
        Log.d(tag, "Login failed: $error")
        if (error != null) {
            snackbar_layout.snackbar(error)
        } else {
            Log.w(tag, "Error message was null")
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
}
