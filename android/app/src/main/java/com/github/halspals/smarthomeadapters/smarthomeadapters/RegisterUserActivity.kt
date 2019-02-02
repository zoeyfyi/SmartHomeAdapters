package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import kotlinx.android.synthetic.main.activity_register_user.*
import org.jetbrains.anko.*
import org.jetbrains.anko.design.snackbar
import org.json.JSONObject

class RegisterUserActivity : AppCompatActivity(), RESTResponseListener {

    private val tag = "RegisterUserActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_user)

        /*
            Set up the click events for the buttons.
         */
        register_button.setOnClickListener { _ ->
            register_button.isEnabled = false
            progressBar.visibility = View.VISIBLE
            registerNewUser()
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
    private fun registerNewUser() {
        progressBar.visibility = View.VISIBLE

        val email = email_input.text.toString()
        val password = password_input.text.toString()
        val confirmedPassword = confirm_password_input.text.toString()

        val inputsOK = checkEmailInput(email)
                && checkPasswordInput(password)
                && checkConfirmPasswordInput(password, confirmedPassword)

        if (inputsOK) {
            RESTRequestTask(this).execute(RESTRequest.REGISTER(email, password))
        } else {
            register_button.isEnabled = true
            progressBar.visibility = View.GONE
        }

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


    override fun handleRESTResponse(responseCode: Int, response: String) {
        Log.d(tag, "Auth response: $responseCode; $response")
        if (responseCode == 200) {
            // TODO do we wanna check that the email in the response matches
            // the one we sent?
            toast("Successfully registered your account")
            startActivity(intentFor<MainActivity>().clearTask().newTask())
        } else {
            val responseJSON = JSONObject(response)
            val errorMsg = responseJSON.getString("error")
            snackbar_layout.snackbar(errorMsg)
            register_button.isEnabled = true
            progressBar.visibility = View.GONE
        }
    }
}
