do language plpgsql
$$
    declare
        user1_uuid constant uuid = '6386B7CA-FFCA-4F7E-A047-A3E46A06A56A';
        user2_uuid constant uuid = '6BBA792E-8320-46B4-84BD-B44122E71AA0';
        user3_uuid constant uuid = '1A7BB22F-B31B-466B-986D-3FD3A3E72475';
        user4_uuid constant uuid = 'EC16AA32-2919-4A09-8194-53975BDE37FC';
        user5_uuid constant uuid = '4F417BA6-B94F-4EEE-A7B1-B6643FBD5926';

        task1_uuid constant uuid = '82D6E6CD-858C-47D9-94DA-4A877591F28C';
        task2_uuid constant uuid = 'C3755CE6-F672-443E-8073-2450CB7D7A85';
        task3_uuid constant uuid = '64BB7CA6-5BD2-4D4A-8E87-7F59AE34F0FB';

        nil_uuid constant uuid = '00000000-0000-0000-0000-000000000000';

    begin
        INSERT INTO app.task_events (task_uuid, event, user_uuid, timestamp)
        VALUES
               (task1_uuid, 'created', user1_uuid, '2022-06-10 10:38:40+00'),
               (task1_uuid, 'approved_by', user2_uuid, '2022-06-10 11:00:00+00'),
               (task1_uuid, 'approved_by', user3_uuid, '2022-06-10 11:05:00+00'),
               (task1_uuid, 'approved_by', user4_uuid, '2022-06-10 12:30:40+00'),
               (task1_uuid, 'signed', nil_uuid, '2022-06-10 12:30:40+00'),
               (task1_uuid, 'sent', nil_uuid, '2022-06-10 12:31:00+00'),
               (task2_uuid, 'created', user1_uuid, '2022-06-10 10:38:40+00'),
               (task2_uuid, 'approved_by', user5_uuid, '2022-06-11 11:00:00+00'),
               (task2_uuid, 'rejected_by', user3_uuid, '2022-06-11 11:05:00+00'),
               (task2_uuid, 'approved_by', user2_uuid, '2022-06-10 12:30:40+00'),
               (task1_uuid, 'created', user3_uuid, '2022-05-09 11:38:41+00'),
               (task3_uuid, 'approved_by', user5_uuid, '2022-05-12 11:00:00+00'),
               (task3_uuid, 'approved_by', user2_uuid, '2022-05-20 11:05:00+00'),
               (task3_uuid, 'approved_by', user4_uuid, '2022-05-26 12:30:40+00'),
               (task3_uuid, 'approved_by', user1_uuid, '2022-06-07 12:30:40+00'),
               (task3_uuid, 'signed', nil_uuid, '2022-06-07 12:31:00+00'),
               (task3_uuid, 'sent', nil_uuid, '2022-06-07 12:31:10+00');
    end;
$$