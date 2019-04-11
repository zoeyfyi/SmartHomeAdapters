package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v4.app.FragmentManager
import android.util.Log
import android.widget.LinearLayout
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.activity_main.*
import net.openid.appauth.AuthState
import net.openid.appauth.AuthorizationService
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.intentFor
import org.jetbrains.anko.newTask
import org.jetbrains.anko.toast

/**
 * The activity encompassing the app's main [Fragment]s.
 */
class MainActivity : AppCompatActivity(), RestApiActivity {

    override val restApiService by lazy {
        RestApiService.new()
    }

    override val authState by lazy {
        readAuthState(this)
    }

    override val authService by lazy {
        AuthorizationService(this)
    }

    override fun makeToast(charSequence: CharSequence) {
        // TODO rewrite this using function referencing
        this.toast(charSequence)
    }

    override val context: Context
        get() = this

    override val snackbarLayout: LinearLayout
        get() = this.snackbar_layout

    override var isInEditMode = false
    override lateinit var robotToEdit: Robot

    private val tag = "MainActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
    }

    override fun onStart() {
        super.onStart()
        // Start the robots fragment by default
        startFragment(RobotsFragment())
    }

    /**
     * Replaces the currently active fragment, if there is any to replace.
     *
     * @param fragment the Fragment to replace the currently active one with.
     * @param addToBackstack if true, fragment will be added to the backstack,
     * otherwise backstack will be dropped
     */
    override fun startFragment(fragment: Fragment, addToBackstack: Boolean, args: Bundle?) {
        Log.d(tag, "[startFragment] Invoked")

        val fManager = supportFragmentManager
        fManager.beginTransaction().run {
            replace(R.id.fragmentContainer, fragment)

            // manually handle the backstack
            if (addToBackstack) {
                // A->B to A->B->C (add to backstack)
                addToBackStack(null)
            } else {
                // A->B->C to A (clear backstack)
                fManager.popBackStack(null, FragmentManager.POP_BACK_STACK_INCLUSIVE)
            }

            commit()
        }

        Log.d(tag, "[startFragment] Committed transaction to fragment")
    }

    internal fun signOut() {
        Log.v(tag, "[signOut] Invoked")
        writeAuthState(this, AuthState())
        startActivity(intentFor<AuthenticationActivity>().clearTask().newTask())
    }
}
