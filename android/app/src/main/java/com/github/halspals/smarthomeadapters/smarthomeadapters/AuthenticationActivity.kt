package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.net.Uri
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import kotlinx.android.synthetic.main.activity_authentication.*
import net.openid.appauth.*
import org.jetbrains.anko.*

class AuthenticationActivity : AppCompatActivity() {

    private val tag = "AuthenticationActivity"

    private val authRequest: AuthorizationRequest by lazy {
        // Otherwise set up Authorization config and get the auth request code
        val authServiceConfig = AuthorizationServiceConfiguration(
                Uri.parse("https://oauth.halspals.co.uk/oauth2/auth"),
                Uri.parse("https://oauth.halspals.co.uk/oauth2/token")
        )
        // Get an authorization code
        AuthorizationRequest.Builder(
                authServiceConfig,
                "b43ce28c-f4e3-412b-8dc5-854a32a0c8db",
                ResponseTypeValues.CODE,
                Uri.parse("http://callback.halspals.co.uk")
        ).build()
    }
    private val authService: AuthorizationService by lazy {
        AuthorizationService(this)
    }
    internal val authState: AuthState by lazy {
        readAuthState(this)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        // Check if there is already an active auth session
        if (authState.isAuthorized) {
            startActivity<MainActivity>()
        }

        login_button.setOnClickListener { _ ->
            startOAuthLogin()
        }
    }

    private fun startOAuthLogin() {
        // Start the authorization webview
        Log.d(tag, "Starting Oauth2 call")
        val authIntent = authService.getAuthorizationRequestIntent(authRequest)
        startActivityForResult(authIntent, 42)
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {

        if (requestCode == 42) {
            Log.d(tag, "Received result from Oauth2 call")
            if (data == null) {
                Log.e(tag, "onActivityResult got null data; aborting")
                return
            }

            // Extract the authorization response
            val resp: AuthorizationResponse? = AuthorizationResponse.fromIntent(data)
            val ex: AuthorizationException? = AuthorizationException.fromIntent(data)
            authState.update(resp, ex)
            if (resp != null) {
                Log.d(tag, "OAuth authorization was successful")
                writeAuthState(this, authState)
                startActivity<MainActivity>()
            } else {
                Log.e(tag, "OAuth failed: got AuthorizationException: $ex")
            }
        } else {
            Log.w(tag, "Received UNEXPECTED activity result")
        }
    }
}
