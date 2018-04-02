from gym.envs.registration import register

register(
    id='physsio-v0',
    entry_point='gym_physsio.envs:PhyssioEnv',
)
