package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigParameter


class ParameterAdapter (private val context: Context, private val parameters: List<ConfigParameter>):  BaseAdapter() {

    override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View {

        if (convertView != null) {
            // We've already been given a view which is set up -- use it
            return convertView
        }

        // Otherwise inflate the appropriate card and set up its fields
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater

        val view = when (parameters[position].type) {
            ConfigParameter.BOOL_TYPE -> {
                // The parameter is a boolean one -- inflate the bool card and set the default state
                val cardView = inflater.inflate(R.layout.view_bool_config_card, parent, false)
                cardView.findViewById<Switch>(R.id.input_switch).isChecked = parameters[position].details.default != 0
                cardView
            }

            ConfigParameter.INT_TYPE -> {
                // The parameter is an integer one -- inflate the int card and set up the seek bar
                val cardView = inflater.inflate(R.layout.view_int_config_card, parent, false)

                val seekBar = cardView.findViewById<SeekBar>(R.id.input_seekbar)
                val max = parameters[position].details.max
                val min = parameters[position].details.min
                seekBar.max = max - min // manually add min later as api <26 doesn't support seekBar.setMin()
                seekBar.progress = parameters[position].details.default - min

                val seekBarValueTextView = cardView.findViewById<TextView>(R.id.seekbar_value_text_view)
                seekBarValueTextView.text = "${parameters[position].details.default}"

                // Set up the change listener for the seekbar
                seekBar.setOnSeekBarChangeListener(object : SeekBar.OnSeekBarChangeListener {
                    override fun onStopTrackingTouch(p0: SeekBar?) {
                        // Not interested
                    }

                    override fun onStartTrackingTouch(p0: SeekBar?) {
                        // Not interested
                    }

                    override fun onProgressChanged(seekBar: SeekBar?, progress: Int, byUser: Boolean) {
                        seekBarValueTextView.text = "${progress + min}"
                    }
                })

                cardView
            }

            else -> {
                TODO("No other types expected")
            }
        }

        // Regardless of the type of card, set up the name and description fields
        view.findViewById<TextView>(R.id.config_name_text_view).text = parameters[position].name
        view.findViewById<TextView>(R.id.config_explanation_text_view).text = parameters[position].description

        return view
    }

    override fun getCount(): Int {
        return parameters.size
    }

    override fun getItem(position: Int): Any {
        return parameters[position]
    }

    override fun getItemId(p0: Int): Long {
        return 0L
    }

}