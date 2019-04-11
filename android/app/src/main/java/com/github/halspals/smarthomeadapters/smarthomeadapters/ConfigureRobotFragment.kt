package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v7.widget.LinearLayoutManager
import android.support.v7.widget.LinearSnapHelper
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigParameter
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.activity_register_robot.*
import kotlinx.android.synthetic.main.fragment_configure_robot.*
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.intentFor
import org.jetbrains.anko.newTask
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

const val RECALIBRATE_FLAG = "Recalibrate"

/**
 * The final screen of the robot registration wizard, where the user configures the robot parameters.
 */
class ConfigureRobotFragment : Fragment() {

    private val fTag = "ConfigureRobotFragment"

    internal val parent by lazy { activity as RegisterRobotActivity }

    internal var numAcksExpected = 0
    internal var numAcksReceived = 0
    internal var numRejectsReceived = 0

    private var getRobotDone = false
    private var getParamsDone = false


    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_configure_robot, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        if (arguments?.getBoolean(RECALIBRATE_FLAG, false) == true) {
            // If we have been instructed to recalibrate, rename the finish button for clarity
            finish_button.setText(R.string.reconfigure_button_text)
        }

        finish_button.setOnClickListener { _ ->
            parent.toast("Finished configuration")
            parent.startActivity(parent.intentFor<MainActivity>().clearTask().newTask())
        }

        cancel_button.setOnClickListener { _ -> parent.finish() }

        parameter_recycler_view.layoutManager = LinearLayoutManager(
                view.context, LinearLayoutManager.HORIZONTAL, false)
        val snapHelper = LinearSnapHelper()
        snapHelper.attachToRecyclerView(parameter_recycler_view)

        progress_bar.visibility = View.VISIBLE
        finish_button.isEnabled = false

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken, _, ex ->
            if (accessToken == null) {
                Log.e(fTag, "Got null access token, ex: $ex")
            } else {
                getConfigParameters(accessToken)
                getRobot(accessToken)

            }
        }
    }

    private fun getRobot(accessToken: String) {
        parent.restApiService
                .getRobot(parent.robotId, accessToken)
                .enqueue(object: Callback<Robot> {

                    override fun onResponse(call: Call<Robot>, response: Response<Robot>) {
                        val robot = response.body()

                        if (response.isSuccessful && robot != null) {
                            Log.d(fTag, "Successfully retrieved $robot")
                            // TODO REMVOE THE BELOW WHEN SERVER FIXED
                            if (robot.robotType == Robot.ROBOT_TYPE_THERMOSTAT) {
                                robot.robotStatus.min = 273
                                robot.robotStatus.max = 373
                                robot.robotStatus.current = 293
                            }
                            setRobotView(robot)
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "Getting the robot was unsuccessful, "
                                    + "error: $error")
                            if (error != null) {
                                parent.snackbar_layout.snackbar(error)
                            }
                        }
                    }

                    override fun onFailure(call: Call<Robot>, t: Throwable) {
                        val error = t.message
                        Log.e(fTag, "fetchRobot() failed: $error")
                        if (error != null) {
                            parent.snackbar_layout.snackbar(error)
                        }
                    }
                })
    }

    /**
     * Inflates a Robot card view into the layout for the robot which is being edited.
     * Also sets the [title_text_view] according to the robot's name.
     */
    private fun setRobotView(robot: Robot) {
        // inflate card view
        val inflater = parent.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val robotView = inflater.inflate(R.layout.view_robot_card, robot_layout.parent as ViewGroup, false)

        // Get internal views
        val robotNickname = robotView.findViewById<TextView>(R.id.robot_nickname_text_view)
        val robotIcon = robotView.findViewById<ImageView>(R.id.robot_image_view)
        val robotCircle = robotView.findViewById<ImageView>(R.id.robot_circle_drawable)
        val robotText = robotView.findViewById<TextView>(R.id.robot_range_text_view)

        // Set up internal views
        robot.updateViews(parent, robotCircle, robotIcon, robotText)
        robot.setViewEvents(parent, robotCircle, robotIcon, robotText)
        robotNickname.text = robot.nickname

        robot_layout.addView(robotView)

        getRobotDone = true
        allowFinishIfAllDone()
    }

    /**
     * Gets the configuration parameters for the robot and updates the [parameter_recycler_view].
     */
    private fun getConfigParameters(accessToken: String) {
        parent.restApiService
                .getConfigParameters(parent.robotId, accessToken)
                .enqueue(object: Callback<List<ConfigParameter>> {

                    override fun onResponse(
                            call: Call<List<ConfigParameter>>,
                            response: Response<List<ConfigParameter>>) {

                        val params = response.body()

                        if (response.isSuccessful && params != null) {
                            Log.v(fTag, "[getConfigParameters] Got params $params")
                            // Set up the grid's adapter to display the configuration
                            // parameters requested
                            parameter_recycler_view.adapter = ParameterAdapter(this@ConfigureRobotFragment, params)
                            getParamsDone = true
                            allowFinishIfAllDone()

                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "[getConfigParameters] got unsuccessful "
                                    + "response or null body; body $params, error: $error")
                            if (error != null) {
                                parent.snackbar_layout.snackbar(error)
                            }
                        }

                    }

                    override fun onFailure(call: Call<List<ConfigParameter>>, t: Throwable) {

                        val error = t.message
                        Log.e(fTag, "[getConfigParameters] FAILED, error: $error")
                        if (error != null) {
                            parent.snackbar_layout.snackbar(error)
                        }

                    }
                })
    }

    internal fun allowFinishIfAllDone() {
        if (getRobotDone && getParamsDone && numAcksExpected == numAcksReceived + numRejectsReceived) {
            progress_bar.visibility = View.GONE
            finish_button.isEnabled = true
        }
    }
}
