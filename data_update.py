from sqlalchemy import create_engine, MetaData
from sqlalchemy.orm import scoped_session, sessionmaker
from sqlalchemy.sql import text

engine = create_engine('postgresql://masterp:@localhost:5432/dev_sms_backup')
# engine = create_engine('postgresql://masterp:pass@localhost:5432/dev_sms_backup')

# metadata = MetaData(bind=engine)
# metadata.reflect()

# result_set = engine.execute('SELECT * FROM students')

# for r in result_set:  
#     print(r[0],r[15])
#     subject_class_query = text("SELECT * FROM subject_classes WHERE is_compulsory = TRUE and class = :student_class;")
#     print(subject_class_query)
#     subject_class = engine.execute(subject_class_query, student_class=r[15])
#     for i in subject_class:
#         print(i)

#         # insert_student_query = text("INSERT into student_subject_classes (student_id, subject_class_id, is_active, created_at, updated_at) VALUES(2, 570, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);")
#         insert_student_query = text("INSERT into student_subject_classes (student_id, subject_class_id, is_active, created_at, updated_at) VALUES(:student_id, :subject_class_id, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);")
#         print(f"The student id is {r[0]} and subject class id is {i[0]}")
#         student_subject_insert = engine.execute(insert_student_query, student_id=r[0], subject_class_id=i[0])


# result_set_ss1 =  engine.execute("SELECT * FROM students WHERE class_text = 'SSS1'")

# ss1_subject_results = engine.execute("SELECT id FROM student_subject_classes WHERE student_id = 21;")

# subject_class_ids = []
# for i in ss1_subject_results:
#     subject_class_ids.append(i[0])
# import pdb; pdb.set_trace()

# for r in result_set_ss1:
#     print(r[0])

    
#     ss1_subject_results = engine.execute("SELECT id FROM student_subject_classes WHERE student_id = 21;")
#     for i in ss1_subject_results:
#         # print(i, r)
#         insert_student_query = text("INSERT into student_subject_classes (student_id, subject_class_id, is_active, created_at, updated_at) VALUES(:student_id, :subject_class_id, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);")
#         print(f"The student id is {r[0]} and subject class id is {i[0]}")
#         student_subject_insert = engine.execute(insert_student_query, student_id=r[0], subject_class_id=i[0])



result_set_to_destroy = engine.execute("SELECT * FROM student_subject_classes ORDER BY student_id")

student_dict = {}

# import pdb; pdb.set_trace()
for i in result_set_to_destroy:
    if student_dict.get(i[4]):
        student_class = engine.execute(f"SELECT class_text FROM students WHERE id={i[4]}").fetchone()
        subject_class = engine.execute(f"SELECT class FROM subject_classes WHERE id={i[5]}").fetchone()
        if student_class[0] != subject_class[0]:
            print(f'The following were duplicated student_id = {i[4]} subject_dict= {i[5]} id = {i[0]} ')
            engine.execute(f"DELETE FROM student_subject_classes WHERE id={i[0]}")
        student_dict[i[4]].update({
            i[5]: i[0]
        })
        

        # if student current class is not the same as subject class delete student subject class record
        # if 
        # if student_dict.get(i[4]).get(i[5]):
        #     print(f'The following were duplicated student_id = {i[4]} subject_dict= {i[5]} id = {i[0]} ')
        #     engine.execute(f"DELETE FROM student_subject_classes WHERE id={i[0]}")
        # student_dict[i[4]].update({
        #     i[5]: i[0]
        # })
    else:
        student_dict[i[4]] = {
            i[5]: i[0]
        }

print(student_dict)