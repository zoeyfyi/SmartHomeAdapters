package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.net.Uri
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import kotlinx.android.synthetic.main.activity_authentication.*
import net.openid.appauth.*
import org.jetbrains.anko.*

/**
 * A screen which interacts with the OAuth server to provide authorization of the user.
 */
class AuthenticationActivity : AppCompatActivity() {

    private val tag = "AuthenticationActivity"

    private val authRequestCode = 42

    private lateinit var authRequest: AuthorizationRequest
    private lateinit var authService: AuthorizationService
    private lateinit var authState: AuthState

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_authentication)

        auth_button.setOnClickListener { _ ->
            startOAuthView()
        }

        val authServiceConfig = AuthorizationServiceConfiguration(
                Uri.parse("https://oauth.halspals.co.uk/oauth2/auth"),
                Uri.parse("https://oauth.halspals.co.uk/oauth2/token")
        )
        authRequest = AuthorizationRequest.Builder(
                authServiceConfig,
                "refresh_test4",
                ResponseTypeValues.CODE,
                Uri.parse("https://callback.halspals.co.uk")
        ).setScope("offline").build()

        authService = AuthorizationService(this)
    }

    override fun onStart() {
        super.onStart()

        authState = readAuthState(this)
        // Check if there is already an active auth session
        if (authState.isAuthorized) {
            startActivity<MainActivity>()
        }
    }

    /**
     * Starts the authorization web view, allowing the user to login or register.
     */
    private fun startOAuthView() {
        Log.d(tag, "Starting Oauth2 call")
        val authIntent = authService.getAuthorizationRequestIntent(authRequest)
        startActivityForResult(authIntent, authRequestCode)
    }

    /**
     * Listens for the result of the OAuth web view.
     */
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

    /**
     * Starts an [TokenRequest] to exchange an authorization code for authorization tokens.
     * Updates the [AuthState] on the device.
     *
     * @param response the [AuthorizationResponse] with the received auth code
     */
    private fun exchangeCodeForTokens(response: AuthorizationResponse) {
        Log.d(tag, "[exchangeCodeForTokens] Starting code-for-tokens exchange")
        authService.performTokenRequest(response.createTokenExchangeRequest())
        { tokenResponse: TokenResponse?, ex: AuthorizationException? ->
            authState.update(tokenResponse, ex)
            writeAuthState(this, authState)
            if (tokenResponse != null) {
                Log.d(tag, "[exchangeCodeForTokens] Exchange successful, moving to main")
                startActivity(intentFor<MainActivity>().newTask())
            } else {
                Log.e(tag, "[exchangeCodeForTokens] Failed, exception $ex")
                // TODO make the user try again etc
            }
        }
    }
}
