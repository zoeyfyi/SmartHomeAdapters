package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import kotlinx.android.synthetic.main.activity_authentication.*

class AuthenticationActivity : AppCompatActivity() {

    private val tag = "AuthenticationActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        // Set up appropriate listeners for button click events
        sign_in_button.setOnClickListener { _ ->
            if (signInUser()) {
                startMainActivity()
            } else {
                Log.w(tag, "Sign-in failed.")
            }
        }

        register_button.setOnClickListener { _ ->
            if (registerNewUser()) {
                startMainActivity()
            } else {
                Log.w(tag, "Registration failed.")
            }
        }

    }

    private fun signInUser(): Boolean {
        // TODO make this authenticate with web server
        return true
    }

    private fun registerNewUser(): Boolean {
        // TODO make this authenticate with the web server, setting up a new user account
        return true
    }

    private fun startMainActivity() {
        val intent = Intent(this, MainActivity::class.java)
        startActivity(intent)
    }
}
