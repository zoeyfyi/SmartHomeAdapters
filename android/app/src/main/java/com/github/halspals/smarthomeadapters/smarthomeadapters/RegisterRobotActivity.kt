package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v4.app.FragmentManager
import android.util.Log
import com.google.zxing.integration.android.IntentIntegrator
import com.google.zxing.integration.android.IntentResult

class RegisterRobotActivity : AppCompatActivity() {

    private val tag = "RegisterRobotActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_robot)

        startFragment(QRFragment())
    }

    /*
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

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {

        val result: IntentResult? = IntentIntegrator.parseActivityResult(requestCode, resultCode, data)
        val robotId: String? = result?.contents

        if (robotId != null) {
            Log.d(tag, "[onActivityResult] Scanned robotId $robotId")

            val registerRobotFragment = RobotSetupFragment()
            val args = Bundle()
            args.putString("robotId", robotId)
            registerRobotFragment.arguments = args

            startFragment(registerRobotFragment, false)
        } else {
            Log.d(tag, "User quit QR scanner early")
        }
    }
}
