package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.widget.LinearLayout
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import net.openid.appauth.AuthState
import net.openid.appauth.AuthorizationService

/**
 * Provides an interface for activities which allow for making requests to the web server.
 *
 * @property restApiService the activity's [RestApiService] instance
 * @property authState the activity's [AuthState] instance for user auth
 * @property authService the activity's [AuthorizationService] instance for user auth
 * @property snackbarLayout the [LinearLayout] in the activity to display snackbars in
 * @property context the [Context] of the activity
 * @property isInEditMode whether the activity is in edit-robot-mode; only applies to [MainActivity]s
 * @property robotToEdit the robot which the activity is wanting to edit, if there is one; only applies to [MainActivity]s
 * @property makeToast the function to use to display toasts in the activity
 * @property startFragment the activity's function for starting new fragments
 */
interface RestApiActivity {
    val restApiService: RestApiService
    val authState: AuthState
    val authService: AuthorizationService
    val snackbarLayout: LinearLayout
    val context: Context
    var isInEditMode: Boolean
    var robotToEdit: Robot
    fun makeToast(charSequence: CharSequence)
    fun startFragment(fragment: Fragment, addToBackstack: Boolean = false, args: Bundle? = null)
}