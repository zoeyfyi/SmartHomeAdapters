package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import kotlinx.android.synthetic.main.activity_register_user.*
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.startActivity
import org.jetbrains.anko.toast

class RegisterUserActivity : AppCompatActivity() {

    private val tag = "RegisterUserActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_user)

        /*
            Set up the click events for the buttons.
         */
        email_image_view.setOnClickListener { _ ->
            // Fire up the user's default email app.
            val intent = Intent(Intent.ACTION_MAIN)
            intent.addCategory(Intent.CATEGORY_APP_EMAIL)
            startActivity(intent)
        }

        send_email_button.setOnClickListener { _ ->
            send_email_button.isEnabled = false
            if (registerNewUser()) {
                switchToVerificationContext()
            } else {
                // TODO the below should also include an error message from the server
                snackbar_layout.snackbar("Registration failed")
                send_email_button.isEnabled = true
            }
        }

        verify_code_button.setOnClickListener { _ ->
            verify_code_button.isEnabled = false
            if (verifyActivationCode()) {
                toast("Successfully registered your account")
                startActivity<MainActivity>()
            } else {
                // TODO the below should also include an error message from the server
                snackbar_layout.snackbar("Verification failed")
                verify_code_button.isEnabled = true
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
                checkEmailInput()
            }
        }

        password_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                checkPasswordInput()
            }
        }

        confirm_password_input.setOnFocusChangeListener { _, hasFocus ->
            if (!hasFocus) {
                checkConfirmPasswordInput()
            }
        }


    }

    /**
     * WIP: Gets the input email and password and registers a new user with the web service.
     * Currently just makes sure that all three inputs are valid.
     *
     * @return whether setting up the new user was successful
     */
    private fun registerNewUser(): Boolean {
        progressBar.visibility = View.VISIBLE
        // TODO read data from input fields and send it to the server, moving on only when
        // receiving a success
        val inputsOK = checkEmailInput() && checkPasswordInput() && checkConfirmPasswordInput()
        progressBar.visibility = View.GONE


        return inputsOK
    }

    /**
     * Hides the views corresponding to registering, instead showing those
     * where the user is asked to verify their email before proceeding.
     */
    private fun switchToVerificationContext() {
        // Hides the previous input fields and buttons and displays
        // new ones as appropriate to verify the user's email
        email_input_layout.visibility = View.GONE
        password_input_layout.visibility = View.GONE
        confirm_password_input_layout.visibility = View.GONE
        send_email_button.visibility = View.GONE
        login_button.visibility = View.GONE
        email_image_view.visibility = View.VISIBLE
        email_sent_textView.visibility = View.VISIBLE
        activation_code_input.visibility = View.VISIBLE
        verify_code_button.visibility = View.VISIBLE
    }

    /**
     * WIP: Gets the verification code from [activation_code_input] and checks with the
     * web service that it looks as expected.
     * Currently just a dummy function; returns true regardless.
     */
    private fun verifyActivationCode(): Boolean {
        progressBar.visibility = View.VISIBLE
        // TODO read data from input fields and send it to the server, moving on only when
        // receiving a success
        progressBar.visibility = View.GONE

        return true
    }

    /**
     * Gets the email from the [email_input] and checks that it is a valid email string.
     * Sets errors on [email_input] as appropriate.
     *
     * @return whether the given email is a valid one
     */
    private fun checkEmailInput(): Boolean {
        val email = email_input.text
        val emailIsValid = (!email.isNullOrEmpty()
                && android.util.Patterns.EMAIL_ADDRESS.matcher(email).matches())

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

    /**
     * Checks that the password input in [confirm_password_input] matches that in [password_input].
     *
     * @return whether the two passwords match
     */
    private fun checkConfirmPasswordInput(): Boolean {
        // Simply make sure that the two passwords given match.

        val confirmedPwIsValid = (confirm_password_input.text.toString()
                == password_input.text.toString())

        if (!confirmedPwIsValid) {
            confirm_password_input.error = "Passwords do not match"
        } else {
            confirm_password_input.error = null
        }

        return confirmedPwIsValid
    }
}
