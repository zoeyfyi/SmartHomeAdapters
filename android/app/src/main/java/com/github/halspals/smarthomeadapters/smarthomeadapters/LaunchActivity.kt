package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.intentFor
import org.jetbrains.anko.newTask

class LaunchActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_launch)

        val savedToken = fetchSavedToken()
        if (savedToken != null && tokenIsValid(savedToken)) {
            startActivity(intentFor<MainActivity>("token" to savedToken).clearTask().newTask())
            overridePendingTransition(android.R.anim.fade_in, android.R.anim.fade_out)
        } else {
            val newToken = tryRefreshAccessToken()
            if (newToken != null) {
                saveToken(newToken)
                startActivity(intentFor<MainActivity>("token" to savedToken).clearTask().newTask())
                overridePendingTransition(android.R.anim.fade_in, android.R.anim.fade_out)
            } else {
                // The current token is invalid, we tried to generate a new one and failed -- the user
                // must thus log in again
                startActivity(intentFor<AuthenticationActivity>().clearTask().newTask())
                overridePendingTransition(android.R.anim.fade_in, android.R.anim.fade_out)
            }
        }
    }

    private fun fetchSavedToken(): String? {
        val prefs = getSharedPreferences(SHARED_PREFERENCES_FILE, Context.MODE_PRIVATE)

        return prefs.getString("token", null)
    }

    private fun saveToken(token: String) {
        val prefs = getSharedPreferences(SHARED_PREFERENCES_FILE, Context.MODE_PRIVATE)
        val editor = prefs.edit()
        editor.putString("token", token)
        editor.apply()
    }

    private fun tokenIsValid(token: String): Boolean {
        // TODO check with the server if this token should still be valid
        return true
    }

    private fun tryRefreshAccessToken(): String? {
        // TODO fetch the refresh token and try to generate a new access token
        return null
    }
}
