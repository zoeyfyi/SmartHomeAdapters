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

    private val authRequestCode = 42

    private val authRequest: AuthorizationRequest by lazy {
        val authServiceConfig = AuthorizationServiceConfiguration(
                Uri.parse("https://oauth.halspals.co.uk/oauth2/auth"),
                Uri.parse("https://oauth.halspals.co.uk/oauth2/token")
        )
        AuthorizationRequest.Builder(
                authServiceConfig,
                "65a0a8b8-9175-4f12-a270-461cb2e8fd85",
                ResponseTypeValues.CODE,
                Uri.parse("https://callback.halspals.co.uk")
        ).setScope("openid").build()
    }

    private val authService: AuthorizationService by lazy {
        AuthorizationService(this)
    }

    private val authState: AuthState by lazy {
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

        register_button.setOnClickListener { _ ->
            startActivity<RegisterUserActivity>()
        }

        if (intent.getBooleanExtra(RegisterUserActivity.FORCE_SIGN_IN, false)) {
            startOAuthLogin()
        }
    }

    private fun startOAuthLogin() {
        // Start the authorization webview
        Log.d(tag, "Starting Oauth2 call")
        val authIntent = authService.getAuthorizationRequestIntent(authRequest)
        startActivityForResult(authIntent, authRequestCode)
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {

        if (requestCode == authRequestCode) {
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
                Log.d(tag, "OAuth authorization was successful: code ${resp.authorizationCode}")
                writeAuthState(this, authState)
                exchangeCodeForTokens(resp)
            } else {
                Log.e(tag, "OAuth failed: got AuthorizationException: $ex")
            }
        } else {
            Log.w(tag, "Received UNEXPECTED activity result with requestCode $requestCode")
        }
    }

    private fun exchangeCodeForTokens(response: AuthorizationResponse) {
        Log.d(tag, "[exchangeCodeForTokens] Starting code-for-tokens exchange")
        authService.performTokenRequest(response.createTokenExchangeRequest())
        { tokenResponse: TokenResponse?, ex: AuthorizationException? ->
            authState.update(tokenResponse, ex)
            writeAuthState(this, authState)
            if (tokenResponse != null) {
                Log.d(tag, "[exchangeCodeForTokens] Exchange successful, moving to main")
                startActivity<MainActivity>()
            } else {
                Log.e(tag, "[exchangeCodeForTokens] Failed, exception $ex")
                // TODO make the user try again etc
            }
        }
    }
}
