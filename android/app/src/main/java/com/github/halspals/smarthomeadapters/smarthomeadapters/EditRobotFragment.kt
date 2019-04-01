package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.activity_main.*
import kotlinx.android.synthetic.main.fragment_edit_robot.*
import okhttp3.ResponseBody
import org.jetbrains.anko.*
import org.jetbrains.anko.design.snackbar
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * A screen which presents the user with various options to edit a robot already added to their account.
 */
class EditRobotFragment : Fragment() {

    private val fTag = "QRFragment"

    private val parent by lazy { activity as MainActivity }

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_edit_robot, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        finish_edit_image_view.setOnClickListener { _ -> parent.startFragment(RobotsFragment()) }

        recalibrate_layout.setOnClickListener { _ -> /*TODO*/ }
        rename_layout.setOnClickListener { _ ->
            parent.startActivity(
                    parent.intentFor<RegisterRobotActivity>(
                            SKIP_TO_SCREEN_FLAG to RENAME_FLAG,
                            ROBOT_ID_FLAG to parent.robotToEdit.id
                    )
            )
        }

        delete_layout.setOnClickListener { _ ->
            parent.alert(
                    "Are you sure you want to delete ${parent.robotToEdit.nickname}?\n" +
                            "This will permanently remove the robot from your account.",
                    "Permanently delete robot") {
                yesButton { deleteRobot(parent.robotToEdit.id) }
                noButton {}
            }.show()
        }

        setRobotViewAndTitle(parent.robotToEdit)
    }

    /**
     * Inflates a Robot card view into the layout for the robot which is being edited.
     * Also sets the [title_text_view] according to the robot's name.
     */
    private fun setRobotViewAndTitle(robot: Robot) {
        // inflate card view
        val inflater = parent.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val robotView = inflater.inflate(R.layout.view_robot_card, robot_layout.parent as ViewGroup, false)

        // get internal views
        val robotNickname = robotView.findViewById<TextView>(R.id.robot_nickname_text_view)
        val robotIcon = robotView.findViewById<ImageView>(R.id.robot_image_view)

        // Set the icon accordingly
        when (robot.robotType) {
            Robot.ROBOT_TYPE_SWITCH -> {
                robotIcon.setImageResource(R.drawable.ic_light_on)
            }

            Robot.ROBOT_TYPE_THERMOSTAT -> {
                robotIcon.setImageResource(R.drawable.basic_accelerator)
            }

            else -> TODO("NO OTHER ROBOT TYPE EXPECTED")

        }

        robotNickname.text = robot.nickname

        robot_layout.addView(robotView)

        title_text_view.text = context?.getString(R.string.erf_title_text, robot.nickname)
    }

    /**
     * Deletes the robot with the matching ID from the user's account.
     *
     * @param robotId the unique ID of the robot to remove
     */
    private fun deleteRobot(robotId: String) {
        parent.authState.performActionWithFreshTokens(parent.authService) {
            accessToken, _, ex ->
            if (accessToken == null) {
                Log.e(fTag, "[deleteRobot] got null access token, ex: $ex")
            } else {

                progress_bar.visibility = View.VISIBLE

                parent.restApiService
                        .deleteRobot(robotId, accessToken)
                        .enqueue(object: Callback<ResponseBody> {

                            override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {
                                progress_bar.visibility = View.GONE

                                if (response.isSuccessful) {
                                    Log.v(fTag, "[deleteRobot] Success")
                                    parent.toast("Deleted robot")
                                    parent.startFragment(RobotsFragment())
                                } else {
                                    val error = RestApiService.extractErrorFromResponse(response)
                                    Log.e(fTag, "[deleteRobot] got unsuccessful "
                                            + "response, error: $error")
                                    if (error != null) {
                                        parent.snackbar_layout.snackbar(error)
                                    }
                                }
                            }

                            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                                progress_bar.visibility = View.GONE
                                val error = t.message
                                Log.e(fTag, "[deleteRobot] FAILED, error: $error")
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }
                            }
                        })
            }
        }
    }

}
