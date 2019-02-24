package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v4.app.FragmentManager
import android.util.Log
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.UseCase
import com.google.zxing.integration.android.IntentIntegrator
import com.google.zxing.integration.android.IntentResult

/**
 * The activity forming the base of the robot registration wizard.
 */
class RegisterRobotActivity : AppCompatActivity() {

    private val tag = "RegisterRobotActivity"

    // Keep track of the id of the robot once it has been scanned/entered
    internal lateinit var robotId: String

    // Record the chosen use case so that its parameters can be set up
    internal lateinit var chosenUseCase: UseCase

    internal val restApiService by lazy {
        RestApiService.new()
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_robot)

        startFragment(QRFragment())
    }

    /*
    * Replaces the currently active fragment, if there is any to replace.
    * Note this IGNORES STATE LOSS.
    *
    * @param fragment the Fragment to replace the currently active one with.
    * @param addToBackstack if true, fragment will be added to the backstack,
    * otherwise backstack will be dropped
    */
    internal fun startFragment(fragment: Fragment, addToBackstack: Boolean = false) {
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
                try {
                    fManager.popBackStack(null, FragmentManager.POP_BACK_STACK_INCLUSIVE)
                } catch (e: IllegalStateException) {
                    Log.w(tag, "[startFragment] Caught IllegalStateException $e")
                }
            }

            commitAllowingStateLoss()
        }

        Log.d(tag, "[startFragment] Committed transaction to fragment")
    }

    /**
     * Handles the result from the ZXING QR Scanner initiated by a child [QRFragment].
     */
    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {

        val result: IntentResult? = IntentIntegrator.parseActivityResult(requestCode, resultCode, data)
        val robotId: String? = result?.contents

        if (robotId != null) {
            Log.d(tag, "[onActivityResult] Scanned robotId $robotId")

            this.robotId = robotId

            val qrFragment = supportFragmentManager.findFragmentById(R.id.fragmentContainer) as QRFragment
            qrFragment.setRobotIdText(robotId)

        } else {
            Log.d(tag, "User quit QR scanner early")
        }
    }
}
