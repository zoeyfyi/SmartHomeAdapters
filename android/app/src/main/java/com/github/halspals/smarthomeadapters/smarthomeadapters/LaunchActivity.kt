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

        val savedToken = fetchSavedAccessToken()
        if (savedToken != null && accessTokenIsValid(savedToken)) {
            startActivity(intentFor<MainActivity>(ACCESS_TOKEN_KEY to savedToken).clearTask().newTask())
            overridePendingTransition(android.R.anim.fade_in, android.R.anim.fade_out)
        } else {
            val newToken = tryRefreshAccessToken()
            if (newToken != null) {
                saveAccessToken(newToken)
                startActivity(intentFor<MainActivity>(ACCESS_TOKEN_KEY to savedToken).clearTask().newTask())
                overridePendingTransition(android.R.anim.fade_in, android.R.anim.fade_out)
            } else {
                // The current token is invalid, we tried to generate a new one and failed -- the user
                // must thus log in again
                startActivity(intentFor<AuthenticationActivity>().clearTask().newTask())
                overridePendingTransition(android.R.anim.fade_in, android.R.anim.fade_out)
            }
        }
    }

    /**
     * Fetches the access token currently stored on the device, if any.
     *
     * @return the token retrieved, or null if there was none
     */
    private fun fetchSavedAccessToken(): String? {
        val prefs = getSharedPreferences(SHARED_PREFERENCES_FILE, Context.MODE_PRIVATE)

        return prefs.getString(ACCESS_TOKEN_KEY, null)
    }

    /**
     * Saves the given access token on the user's device.
     *
     * @param token the access token to store
     */
    private fun saveAccessToken(token: String) {
        val prefs = getSharedPreferences(SHARED_PREFERENCES_FILE, Context.MODE_PRIVATE)
        val editor = prefs.edit()
        editor.putString(ACCESS_TOKEN_KEY, token)
        editor.apply()
    }

    /**
     * WIP: Checks the validity of the access token given.
     * This should communicate with the server to establish the token's validity.
     */
    private fun accessTokenIsValid(token: String): Boolean {
        // TODO check with the server if this token should still be valid
        return true
    }

    /**
     * WIP: Fetches the refresh token and tries to generate a new access token.
     * Currently only returns null, signifying failure, as there are no refresh tokens
     * implemented.
     */
    private fun tryRefreshAccessToken(): String? {
        // TODO fetch the refresh token and try to generate a new access token
        return null
    }
}
