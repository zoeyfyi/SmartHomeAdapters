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
import net.openid.appauth.AuthorizationService

const val SKIP_TO_SCREEN_FLAG = "SkipTo"
const val ROBOT_ID_FLAG = "passedRobotId"

/**
 * The activity forming the base of the robot registration wizard.
 */
class RegisterRobotActivity : AppCompatActivity() {

    private val tag = "RegisterRobotActivity"

    // Keep track of the id of the robot once it has been scanned/entered
    internal lateinit var robotId: String
    internal lateinit var robotNickname: String
    // Record the chosen use case so that its parameters can be set up
    internal lateinit var chosenUseCase: UseCase

    internal val restApiService by lazy {
        RestApiService.new()
    }

    internal val authState by lazy {
        readAuthState(this)
    }

    internal val authService by lazy {
        AuthorizationService(this)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_robot)
    }

    override fun onStart() {
        super.onStart()

        val goToScreen = if (intent?.hasExtra(SKIP_TO_SCREEN_FLAG) == true) {
            intent?.getStringExtra(SKIP_TO_SCREEN_FLAG)
        } else {
            null
        }

        val givenRobotId = if (intent?.hasExtra(ROBOT_ID_FLAG) == true) {
            intent?.getStringExtra(ROBOT_ID_FLAG)
        } else {
            null
        }

        when(goToScreen) {
            RENAME_FLAG -> {
                // The user wants to rename an existing robot.
                // Make sure we have been told which robot to rename
                if (givenRobotId == null) {
                    Log.e(tag,  "User wanted to move straight to NickNameFragment but no " +
                            "robot ID was given")
                } else {
                    // Start the NicknameFragment and let it know we want to rename
                    Log.v(tag, "Starting NicknameFragment for renaming $givenRobotId")
                    robotId = givenRobotId
                    val args = Bundle()
                    args.putBoolean(RENAME_FLAG, true)
                    startFragment(fragment = NicknameFragment(), args = args)
                }
            }

            null -> {
                // Default case; the robot registration wizard should run in full.
                // Start its first screen
                startFragment(QRFragment())
            }
        }
    }

    /**
     * Replaces the currently active fragment, if there is any to replace.
     * Note this IGNORES STATE LOSS.
     *
     * @param fragment the Fragment to replace the currently active one with.
     * @param addToBackstack if true, fragment will be added to the backstack,
     * otherwise backstack will be dropped
     * @param args the arguments to attach to the fragment, if any
     */
    internal fun startFragment(
            fragment: Fragment, addToBackstack: Boolean = false, args: Bundle? = null) {

        Log.d(tag, "[startFragment] Invoked")

        fragment.arguments = args

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
            startFragment(NicknameFragment())

        } else {
            Log.d(tag, "User quit QR scanner early")
        }
    }
}
