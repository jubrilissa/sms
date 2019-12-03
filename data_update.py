from sqlalchemy import create_engine, MetaData
from sqlalchemy.orm import scoped_session, sessionmaker
from sqlalchemy.sql import text

engine = create_engine('postgresql://masterp:@localhost:5432/dev_sms_backup')
# engine = create_engine('postgresql://masterp:pass@localhost:5432/dev_sms_backup')

# metadata = MetaData(bind=engine)
# metadata.reflect()

result_set = engine.execute('SELECT * FROM students')

for r in result_set:  
    print(r[0],r[15])
    subject_class_query = text("SELECT * FROM subject_classes WHERE is_compulsory = TRUE and class = :student_class;")
    print(subject_class_query)
    subject_class = engine.execute(subject_class_query, student_class=r[15])
    for i in subject_class:
        print(i)

        # insert_student_query = text("INSERT into student_subject_classes (student_id, subject_class_id, is_active, created_at, updated_at) VALUES(2, 570, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);")
        insert_student_query = text("INSERT into student_subject_classes (student_id, subject_class_id, is_active, created_at, updated_at) VALUES(:student_id, :subject_class_id, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);")
        print(f"The student id is {r[0]} and subject class id is {i[0]}")
        student_subject_insert = engine.execute(insert_student_query, student_id=r[0], subject_class_id=i[0])