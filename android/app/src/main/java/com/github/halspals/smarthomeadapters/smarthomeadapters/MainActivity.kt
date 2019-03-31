package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v4.app.FragmentManager
import android.util.Log
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import net.openid.appauth.AuthorizationService

/**
 * The activity encompassing the app's main [Fragment]s.
 */
class MainActivity : AppCompatActivity() {

    internal val restApiService by lazy {
        RestApiService.new()
    }

    internal val authState by lazy {
        readAuthState(this)
    }

    internal val authService by lazy {
        AuthorizationService(this)
    }

    internal var isInEditMode = false
    internal lateinit var robotToEdit: Robot

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
    fun startFragment(fragment: Fragment, addToBackstack: Boolean = false) {
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
}
