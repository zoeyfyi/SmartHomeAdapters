package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import android.util.Log
import android.view.View
import kotlinx.android.synthetic.main.activity_authentication.*
import org.jetbrains.anko.*
import org.jetbrains.anko.design.snackbar
import org.json.JSONObject
import java.net.HttpURLConnection
import java.security.KeyStore
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec

class AuthenticationActivity : AppCompatActivity(), RESTResponseListener {

    private val tag = "AuthenticationActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        /*
            Set up appropriate listeners for button click events
         */
        sign_in_button.setOnClickListener { _ ->
            sign_in_button.isEnabled = false

            signInUser()
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
    private fun signInUser() {
        login_progress_bar.visibility = View.VISIBLE

        val email = email_input.text.toString()
        val password = password_input.text.toString()
        val inputsOK = checkEmailInput(email) && checkPasswordInput(password)

        if (inputsOK) {
            RESTRequestTask(this).execute(RESTRequest.LOGIN(email, password))
        } else {
            login_progress_bar.visibility = View.GONE
            snackbar_layout.snackbar("Failed to sign you in!")
            Log.w(tag, "Sign-in failed.")
            sign_in_button.isEnabled = true
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

    override fun handleRESTResponse(responseCode: Int, response: String, requestType: String) {

        if (requestType != RESTRequest.LOGIN_TYPE) {
            // Only expect to hear back from login events
            return
        }

        val responseJSON = JSONObject(response)
        if (responseCode < HttpURLConnection.HTTP_BAD_REQUEST) {
            // Save the token which we received securely
            val token = responseJSON.getString("token")
            if (token == null) {
                Log.e(tag, "Login was successful but did not receive a token")
                return
            }


            toast("Signed in")
            Log.d(tag, "Starting MainActivity")
            startActivity(intentFor<MainActivity>("token" to token).clearTask().newTask())
        } else {
            val errorMsg = responseJSON.getString("error")
            snackbar_layout.snackbar(errorMsg)
            login_progress_bar.visibility = View.GONE
            sign_in_button.isEnabled = true
        }
    }
}
