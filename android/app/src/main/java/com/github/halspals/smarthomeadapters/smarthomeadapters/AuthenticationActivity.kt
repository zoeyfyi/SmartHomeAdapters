package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import kotlinx.android.synthetic.main.activity_authentication.*
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.startActivity
import org.jetbrains.anko.toast

class AuthenticationActivity : AppCompatActivity() {

    private val tag = "AuthenticationActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        /*
            Set up appropriate listeners for button click events
         */
        sign_in_button.setOnClickListener { _ ->
            sign_in_button.isEnabled = false

            if (signInUser()) {
                toast("Signed in")
                Log.d(tag, "Starting MainActivity")
                startActivity<MainActivity>()
            } else {
                // TODO this should use the error message received by the server
                snackbar_layout.snackbar("Failed to sign you in!")
                Log.w(tag, "Sign-in failed.")
                sign_in_button.isEnabled = true
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
                checkEmailInput()
            }
        }
        password_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                // The password input has dropped focus; check its validity.
                checkPasswordInput()
            }
        }
    }

    /**
     * WIP: Get the email and password given, validate them, and then
     * try to authenticate the user.
     *
     * @return whether the inputs are valid and the authentication was successful
     */
    private fun signInUser(): Boolean {
        login_progress_bar.visibility = View.VISIBLE
        // TODO make this authenticate with web server
        val signInOK = checkEmailInput() && checkPasswordInput()
        login_progress_bar.visibility = View.GONE
        return signInOK
    }

    /**
     * Gets the email from the [email_input] and checks that it is a valid email string.
     * Sets errors on [email_input] as appropriate.
     *
     * @return whether the given email is a valid one
     */
    private fun checkEmailInput(): Boolean {
        val email = email_input.text
        val emailIsValid = !email.isNullOrEmpty()
                           && android.util.Patterns.EMAIL_ADDRESS.matcher(email).matches()

        if (!emailIsValid) {
            email_input.error = "Invalid email address"
        } else {
            email_input.error = null
        }

        return emailIsValid
    }


    /**
     * WIP: Gets the password from [password_input] and checks that it is valid.
     * Currently only checks that it is not null or empty.
     * Sets errors on [password_input] as appropriate.
     *
     * @return whether the password is a valid one
     */
    private fun checkPasswordInput(): Boolean {
        // TODO we can add more checks here like password length etc
        val pwIsValid = !password_input.text.isNullOrEmpty()

        if (!pwIsValid) {
            password_input.error = "Invalid password"
        } else {
            password_input.error = null
        }

        return pwIsValid
    }
}
