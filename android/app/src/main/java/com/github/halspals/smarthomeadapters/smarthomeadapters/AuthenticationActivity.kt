package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import kotlinx.android.synthetic.main.activity_authentication.*
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast

class AuthenticationActivity : AppCompatActivity(), ButtonUpdater {

    private val tag = "AuthenticationActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        // Set up appropriate listeners for button click events
        sign_in_button.setOnClickListener { _ ->
            sign_in_button.isEnabled = false

            if (signInUser()) {
                toast("Signed in")
                startMainActivity()
            } else {
                // TODO this should use the error message received by the server
                snackbar_layout.snackbar("Failed to sign you in!")
                Log.w(tag, "Sign-in failed.")
                sign_in_button.isEnabled = true
            }
        }

        register_button.setOnClickListener { _ ->
            startRegisterUserActivity()
        }

        // Add text listeners to enable and disable the sign in button according to what the user
        // has entered
        val nonEmptyTextWatcher = NonEmptyTextWatcher(this, sign_in_button)
        email_input.addTextChangedListener(nonEmptyTextWatcher)
        password_input.addTextChangedListener(nonEmptyTextWatcher)

    }

    private fun signInUser(): Boolean {
        login_progress_bar.visibility = View.VISIBLE
        // TODO make this authenticate with web server
        login_progress_bar.visibility = View.GONE
        return true
    }


    private fun startMainActivity() {
        Log.d(tag, "Starting MainActivity")
        val intent = Intent(this, MainActivity::class.java)
        startActivity(intent)
    }

    private fun startRegisterUserActivity() {
        Log.d(tag, "Starting RegisterUserActivity")
        val intent = Intent(this, RegisterUserActivity::class.java)
        startActivity(intent)
    }

    override fun updateButton(button: Button) {
        // Enables the button iff both the email and password input are non-empty (and non-null)
        val passwordEmpty = password_input.text.isNullOrEmpty()
        val emailEmpty = password_input.text.isNullOrEmpty()

        button.isEnabled = !emailEmpty && !passwordEmpty
    }
}
