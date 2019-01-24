package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import kotlinx.android.synthetic.main.activity_register_user.*
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast

class RegisterUserActivity : AppCompatActivity(), ButtonUpdater {

    private val tag = "RegisterUserActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_user)

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
                startMainActivity()
            } else {
                // TODO the below should also include an error message from the server
                snackbar_layout.snackbar("Verification failed")
                verify_code_button.isEnabled = true
            }
        }

        // Add text watchers to enable and disable the send button according to what the user
        // has entered
        val sendButtonTextWatcher = NonEmptyTextWatcher(this, send_email_button)
        email_input.addTextChangedListener(sendButtonTextWatcher)
        password_input.addTextChangedListener(sendButtonTextWatcher)
        confirm_password_input.addTextChangedListener(sendButtonTextWatcher)

        // Similarly add a text watcher to the verification code input field
        // to enable/disable the verify button
        val verifyButtonTextWatcher = NonEmptyTextWatcher(this, verify_code_button)
        activation_code_input.addTextChangedListener(verifyButtonTextWatcher)
    }

    private fun registerNewUser(): Boolean {
        progressBar.visibility = View.VISIBLE
        // TODO read data from input fields and send it to the server, moving on only when
        // receiving a success
        progressBar.visibility = View.GONE


        return true
    }

    private fun switchToVerificationContext() {
        // Hides the previous input fields and buttons and displays
        // new ones as appropriate to verify the user's email
        email_input_layout.visibility = View.GONE
        password_input_layout.visibility = View.GONE
        confirm_password_input_layout.visibility = View.GONE
        send_email_button.visibility = View.GONE
        email_image_view.visibility = View.VISIBLE
        email_sent_textView.visibility = View.VISIBLE
        activation_code_input.visibility = View.VISIBLE
        verify_code_button.visibility = View.VISIBLE
    }

    private fun verifyActivationCode(): Boolean {
        progressBar.visibility = View.VISIBLE
        // TODO read data from input fields and send it to the server, moving on only when
        // receiving a success
        progressBar.visibility = View.GONE

        return true
    }

    private fun startMainActivity() {
        val intent = Intent(this, MainActivity::class.java)
        startActivity(intent)
    }

    override fun updateButton(button: Button) {
        if (button == send_email_button) {
            // Enables the button iff:
            // * all input fields are non-null,
            // * all input fields are non-empty, and
            // * the passwords given match.
            val emailEmpty = email_input.text.isNullOrEmpty()
            val password = password_input.text
            val passwordEmpty = password_input.text.isNullOrEmpty()
            val confirmedPassword = confirm_password_input.text
            val confirmedPasswordEmpty = confirmedPassword.isNullOrEmpty()

            button.isEnabled = (!emailEmpty
                    && !passwordEmpty
                    && !confirmedPasswordEmpty
                    && password.toString() == confirmedPassword.toString())
        } else if (button == verify_code_button) {
            button.isEnabled = !activation_code_input.text.isNullOrEmpty()
        } else {
            Log.w(tag, "[updateButton] Unexpected button encountered.")
        }

    }
}
