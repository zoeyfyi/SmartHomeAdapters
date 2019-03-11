package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.util.Log
import net.openid.appauth.AuthState

private const val TAG = "AuthStateManager"

/**
 * Fetches a saved [AuthState] if there is one, otherwise makes a new one.
 */
internal fun readAuthState(context: Context): AuthState {
    val prefs = context.getSharedPreferences("auth", Context.MODE_PRIVATE)
    val stateJson = prefs.getString("stateJson", null)
    Log.d(TAG, "[readAuthState] Got stateJson $stateJson")
    return if (stateJson != null) {
       AuthState.jsonDeserialize(stateJson)
    } else {
        AuthState()
    }
}

/**
 * Writes the given [AuthState] to the shared preferences.
 */
internal fun writeAuthState(context: Context, state: AuthState) {
    val prefs = context.getSharedPreferences("auth", Context.MODE_PRIVATE)
    Log.d(TAG, "[writeAuthState] Writing auth state ${state.jsonSerializeString()}")
    prefs.edit().putString("stateJson", state.jsonSerializeString()).apply()
}